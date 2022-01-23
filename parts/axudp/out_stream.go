package axudp

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"github.com/rs/zerolog/log"
	"sync"
	"sync/atomic"
	"time"
)

var retryTTL = time.Millisecond * 500

type OutStream struct {
	optional              *outOptionalStream
	optionalConsistently  *outOptionalStream
	mandatory             *outMandatoryStream
	mandatoryConsistently *outMandatoryConsistentlyStream
}

func NewOutStream(outChan chan *pproto.Packet) *OutStream {
	return &OutStream{
		optional:              newOutOptionalStream(outChan, pproto.PacketMode_PM_OPTIONAL),
		optionalConsistently:  newOutOptionalStream(outChan, pproto.PacketMode_PM_OPTIONAL_CONSISTENTLY),
		mandatory:             newOutMandatoryStream(outChan),
		mandatoryConsistently: newOutMandatoryConsistentlyStream(outChan),
	}
}

func (os *OutStream) send(payload []byte, mode pproto.PacketMode) {
	switch mode {
	case pproto.PacketMode_PM_OPTIONAL:
		os.optional.hold(payload)
	case pproto.PacketMode_PM_OPTIONAL_CONSISTENTLY:
		os.optionalConsistently.hold(payload)
	case pproto.PacketMode_PM_MANDATORY:
		os.mandatory.hold(payload)
	case pproto.PacketMode_PM_MANDATORY_CONSISTENTLY:
		os.mandatoryConsistently.hold(payload)
	}
}

func (os *OutStream) response(pck *pproto.Packet) {
	switch pck.Mode {
	case pproto.PacketMode_PM_MANDATORY:
		os.mandatory.response(pck)
	case pproto.PacketMode_PM_MANDATORY_CONSISTENTLY:
		os.mandatoryConsistently.response(pck)
	}
}

type outStreamBase struct {
	mtu  int
	S    chan *pproto.Packet
	Mode pproto.PacketMode
}

func (os *outStreamBase) setMtu(mtu int) {
	os.mtu = mtu
}

func (os *outStreamBase) send(packet *pproto.Packet) {
	log.Trace().Hex("data", packet.Payload).Msg("send")
	packet.Mode = os.Mode
	os.S <- packet
}

type outOptionalStream struct {
	outStreamBase
	lastId uint64
}

func newOutOptionalStream(outChan chan *pproto.Packet, mode pproto.PacketMode) *outOptionalStream {
	return &outOptionalStream{
		outStreamBase: outStreamBase{
			mtu:  400,
			S:    outChan,
			Mode: mode,
		},
		lastId: 0,
	}
}

func (os *outOptionalStream) hold(payload []byte) {
	chunks := getChunk(payload, os.mtu)
	parts := len(chunks)
	id := atomic.AddUint64(&os.lastId, 1)
	for part, chunk := range chunks {
		pck := &pproto.Packet{
			Id:         id,
			PartsCount: uint32(parts),
			Payload:    chunk,
			Parts:      booleansToBytes(indexInBoolean(parts, part)),
			Tag:        pproto.PacketTag_PT_PAYLOAD,
			Mtu:        uint32(os.mtu),
		}
		os.send(pck)
	}
}

type outMandatoryStream struct {
	outStreamBase
	lastId uint64
	out    map[uint64]*outMandatoryHolder
	lock   sync.RWMutex
}

func newOutMandatoryStream(outChan chan *pproto.Packet) *outMandatoryStream {
	return &outMandatoryStream{
		outStreamBase: outStreamBase{
			mtu:  400,
			S:    outChan,
			Mode: pproto.PacketMode_PM_MANDATORY,
		},
		out:  map[uint64]*outMandatoryHolder{},
		lock: sync.RWMutex{},
	}
}

func (os *outMandatoryStream) hold(payload []byte) {
	h := newOutMandatoryHolder(payload, atomic.AddUint64(&os.lastId, 1), os.mtu, os.resend)
	os.lock.Lock()
	os.out[h.id] = h
	os.lock.Unlock()
	os.resend(h)
}
func (os *outMandatoryStream) resend(h *outMandatoryHolder) {
	for _, pck := range h.getPackets() {
		go func(p *pproto.Packet) {
			os.send(p)
		}(pck)
	}
}
func (os *outMandatoryStream) response(packet *pproto.Packet) {
	os.lock.RLock()
	h, ok := os.out[packet.Id]
	os.lock.RUnlock()
	if !ok {
		return
	}
	h.response(packet)
	if h.done {
		h.close <- true
		os.lock.Lock()
		delete(os.out, h.id)
		os.lock.Unlock()
	}
}

type outMandatoryHolder struct {
	id         uint64
	parts      []bool
	partsCount uint32
	done       bool
	payloads   [][]byte
	tag        pproto.PacketTag
	mtu        int
	retryTimer *time.Timer
	close      chan bool
}

func newOutMandatoryHolder(payload []byte, id uint64, mtu int, resend func(pck *outMandatoryHolder)) *outMandatoryHolder {
	chunk := getChunk(payload, mtu)
	parts := len(chunk)
	res := &outMandatoryHolder{
		id:         id,
		parts:      make([]bool, parts),
		partsCount: uint32(parts),
		payloads:   chunk,
		tag:        pproto.PacketTag_PT_PAYLOAD,
		mtu:        mtu,
		retryTimer: time.NewTimer(retryTTL),
		close:      make(chan bool, 1),
	}
	go func() {
	MainLoop:
		for {
			select {
			case <-res.retryTimer.C:
				res.retryTimer.Reset(retryTTL)
				resend(res)
				log.Debug().Uint64("id", res.id).Msg("resend")
			case <-res.close:
				log.Debug().Uint64("id", res.id).Msg("close")
				break MainLoop
			}
		}
	}()
	return res
}
func (oh *outMandatoryHolder) response(packet *pproto.Packet) {
	mergeBoolLists(&oh.parts, bytesToBooleans(packet.Parts))
	if allBooleans(oh.parts) {
		oh.done = true
		oh.retryTimer.Reset(retryTTL)
	}
}
func (oh *outMandatoryHolder) getPackets() []*pproto.Packet {
	var res []*pproto.Packet
	for i, part := range oh.parts {
		if part {
			continue
		}
		p := &pproto.Packet{
			Id:         oh.id,
			PartsCount: oh.partsCount,
			Payload:    oh.payloads[i],
			Parts:      booleansToBytes(indexInBoolean(int(oh.partsCount), i)),
			Tag:        oh.tag,
		}
		res = append(res, p)
	}
	return res
}

type outMandatoryConsistentlyStream struct {
	outStreamBase
	lastId uint64
	h      *outMandatoryHolder
	out    [][]byte
	lock   sync.RWMutex
}

func newOutMandatoryConsistentlyStream(outChan chan *pproto.Packet) *outMandatoryConsistentlyStream {
	return &outMandatoryConsistentlyStream{
		outStreamBase: outStreamBase{
			mtu:  400,
			S:    outChan,
			Mode: pproto.PacketMode_PM_MANDATORY_CONSISTENTLY,
		},
		h:    nil,
		out:  [][]byte{},
		lock: sync.RWMutex{},
	}
}

func (os *outMandatoryConsistentlyStream) next() {
	if os.h != nil || len(os.out) == 0 {
		return
	}
	payload := os.out[0]
	os.out = os.out[1:]
	os.h = newOutMandatoryHolder(payload, atomic.AddUint64(&os.lastId, 1), os.mtu, os.resend)
	os.resend(os.h)
}
func (os *outMandatoryConsistentlyStream) resend(h *outMandatoryHolder) {
	for _, pck := range h.getPackets() {
		go func(p *pproto.Packet) {
			os.send(p)
		}(pck)
	}
}
func (os *outMandatoryConsistentlyStream) response(packet *pproto.Packet) {
	os.lock.Lock()
	defer os.lock.Unlock()
	if os.h == nil {
		os.next()
		return
	}
	os.h.response(packet)
	if os.h.done {
		os.h.close <- true
		os.h = nil
		os.next()
	}
}
func (os *outMandatoryConsistentlyStream) hold(payload []byte) {
	os.lock.Lock()
	defer os.lock.Unlock()
	os.out = append(os.out, payload)
	os.next()
}
