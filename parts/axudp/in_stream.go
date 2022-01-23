package axudp

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

var dropOptionalPacketTTL = time.Millisecond * 5000 // Сколько ждем неполученные пакеты
var waitDuplicatePacket = time.Millisecond * 30000  // Ждем чтоб не повторились паакеты

type InStream struct {
	optional              *inOptionalStream
	optionalConsistently  *inOptionalConsistentlyStream
	mandatory             *inMandatoryStream
	mandatoryConsistently *inMandatoryStream
}

func NewInStream(respChan chan *pproto.Packet, serviceChannel chan []byte) *InStream {
	return &InStream{
		optional:              newInOptionalStream(respChan, serviceChannel),
		optionalConsistently:  newInOptionalConsistentlyStream(respChan, serviceChannel),
		mandatory:             newInMandatoryStream(respChan, pproto.PacketMode_PM_MANDATORY, serviceChannel),
		mandatoryConsistently: newInMandatoryStream(respChan, pproto.PacketMode_PM_MANDATORY_CONSISTENTLY, serviceChannel),
	}
}

func (is *InStream) receive(pck *pproto.Packet) {
	switch pck.Mode {
	case pproto.PacketMode_PM_OPTIONAL:
		is.optional.receive(pck)
	case pproto.PacketMode_PM_OPTIONAL_CONSISTENTLY:
		is.optionalConsistently.receive(pck)
	case pproto.PacketMode_PM_MANDATORY:
		is.mandatory.receive(pck)
	case pproto.PacketMode_PM_MANDATORY_CONSISTENTLY:
		is.mandatoryConsistently.receive(pck)
	}
}

type inOptionalStream struct {
	in   map[uint64]*inOptionalHolder
	C    chan []byte // Payload to service
	R    chan *pproto.Packet
	mode pproto.PacketMode
	lock sync.Mutex
}

func newInOptionalStream(respChan chan *pproto.Packet, serviceChannel chan []byte) *inOptionalStream {
	return &inOptionalStream{
		C:    serviceChannel,
		R:    respChan,
		mode: pproto.PacketMode_PM_OPTIONAL,
		in:   map[uint64]*inOptionalHolder{},
		lock: sync.Mutex{},
	}
}

func (is *inOptionalStream) receive(pck *pproto.Packet) {
	is.lock.Lock()
	defer is.lock.Unlock()
	h, ok := is.in[pck.Id]
	if !ok {
		h = newInOptionalHolder(pck, is.done)
		is.in[pck.Id] = h
	}

	h.hold(pck)
	go func() {
		is.R <- &pproto.Packet{
			Id:         pck.Id,
			Tag:        pproto.PacketTag_PT_DONE,
			Mode:       pck.Mode,
			Parts:      booleansToBytes(h.parts),
			PartsCount: h.partsCount,
			Mtu:        h.mtu,
		}
	}()
}
func (is *inOptionalStream) done(holder *inOptionalHolder, del bool) {
	is.lock.Lock()
	defer is.lock.Unlock()
	if del {
		delete(is.in, holder.id)
	} else {
		is.C <- holder.getPayload()
	}
}

type inOptionalHolder struct {
	id         uint64
	parts      []bool
	payload    [][]byte
	partsCount uint32
	done       bool
	mtu        uint32
	tag        pproto.PacketTag
	closeTimer *time.Timer
	close      chan bool
	complete   chan bool
}

func newInOptionalHolder(pck *pproto.Packet, done func(holder *inOptionalHolder, del bool)) *inOptionalHolder {
	res := &inOptionalHolder{
		id:         pck.Id,
		parts:      make([]bool, pck.PartsCount),
		payload:    make([][]byte, pck.PartsCount),
		partsCount: pck.PartsCount,
		done:       false,
		mtu:        pck.Mtu,
		tag:        pck.Tag,
		closeTimer: time.NewTimer(dropOptionalPacketTTL),
		close:      make(chan bool, 1),
	}
	go func() {
	MainLoop:
		for {
			select {
			case <-res.closeTimer.C:
				break MainLoop
			case b := <-res.close:
				if b {
					done(res, false)
					res.closeTimer.Reset(waitDuplicatePacket)
				} else {
					break MainLoop
				}
			}
		}
		done(res, true)
	}()
	return res
}

func (ic *inOptionalHolder) stop() {
	ic.close <- false
}
func (ic *inOptionalHolder) getPayload() []byte {
	return bytesJoin(ic.payload...)
}
func (ic *inOptionalHolder) hold(pck *pproto.Packet) {
	if ic.done {
		return
	}
	index := booleanIndex(bytesToBooleans(pck.Parts)[:pck.PartsCount])
	if index >= 0 && !ic.parts[index] {
		ic.parts[index] = true
		ic.payload[index] = pck.Payload
		ic.closeTimer.Reset(dropOptionalPacketTTL)
		if allBooleans(ic.parts) {
			ic.done = true
			ic.close <- true
		}
		log.Debug().Bools("parts", ic.parts).Uint64("id", ic.id).Bool("done", ic.done).Msg("hold")
	}
}

type inOptionalConsistentlyStream struct {
	in     *inOptionalConsistentlyHolder
	lastId uint64
	C      chan []byte // Payload to service
	R      chan *pproto.Packet
	mode   pproto.PacketMode
	lock   sync.Mutex
}

func newInOptionalConsistentlyStream(respChan chan *pproto.Packet, serviceChannel chan []byte) *inOptionalConsistentlyStream {
	res := &inOptionalConsistentlyStream{
		mode: pproto.PacketMode_PM_OPTIONAL_CONSISTENTLY,
		lock: sync.Mutex{},
		R:    respChan,
		C:    serviceChannel,
	}
	return res
}

func (is *inOptionalConsistentlyStream) receive(pck *pproto.Packet) {

	is.lock.Lock()
	defer is.lock.Unlock()
	if pck.Id < is.lastId { //DROP
		return
	}
	if is.in == nil || is.lastId < pck.Id {
		is.lastId = pck.Id
		is.in = newInOptionalConsistentlyHolder(pck)
	}
	if is.lastId == pck.Id {
		is.in.hold(pck)
		if is.in.done {
			if !is.in.complete {
				is.C <- is.in.getPayload()
				is.in.complete = true
			}
		}
	}
}

type inOptionalConsistentlyHolder struct {
	id         uint64
	parts      []bool
	payload    [][]byte
	partsCount uint32
	done       bool
	complete   bool
	mtu        uint32
	tag        pproto.PacketTag
}

func newInOptionalConsistentlyHolder(pck *pproto.Packet) *inOptionalConsistentlyHolder {
	res := &inOptionalConsistentlyHolder{
		id:         pck.Id,
		parts:      make([]bool, pck.PartsCount),
		payload:    make([][]byte, pck.PartsCount),
		partsCount: pck.PartsCount,
		done:       false,
		mtu:        pck.Mtu,
		tag:        pck.Tag,
	}
	return res
}

func (ic *inOptionalConsistentlyHolder) getPayload() []byte {
	return bytesJoin(ic.payload...)
}
func (ic *inOptionalConsistentlyHolder) hold(pck *pproto.Packet) {
	if ic.done {
		return
	}
	index := booleanIndex(bytesToBooleans(pck.Parts)[:pck.PartsCount])
	if index >= 0 && !ic.parts[index] {
		ic.parts[index] = true
		ic.payload[index] = pck.Payload
		if allBooleans(ic.parts) {
			ic.done = true
		}
	}
}

type inMandatoryStream struct {
	in   map[uint64]*inOptionalHolder
	C    chan []byte // Payload to service
	R    chan *pproto.Packet
	mode pproto.PacketMode
	lock sync.Mutex
}

func newInMandatoryStream(respChan chan *pproto.Packet, mode pproto.PacketMode, serviceChannel chan []byte) *inMandatoryStream {
	res := &inMandatoryStream{
		mode: mode,
		in:   map[uint64]*inOptionalHolder{},
		lock: sync.Mutex{},
		R:    respChan,
		C:    serviceChannel,
	}
	return res
}

func (is *inMandatoryStream) receive(pck *pproto.Packet) {
	is.lock.Lock()
	defer is.lock.Unlock()
	h, ok := is.in[pck.Id]
	if !ok {
		h = newInOptionalHolder(pck, is.done)
		is.in[pck.Id] = h
	}
	h.hold(pck)
	is.R <- &pproto.Packet{
		Id:         pck.Id,
		Tag:        pproto.PacketTag_PT_DONE,
		Mode:       pck.Mode,
		Parts:      booleansToBytes(h.parts),
		PartsCount: h.partsCount,
		Mtu:        h.mtu,
	}
}

func (is *inMandatoryStream) done(holder *inOptionalHolder, del bool) {
	is.lock.Lock()
	defer is.lock.Unlock()
	if del {
		delete(is.in, holder.id)
	} else {
		is.C <- holder.getPayload()
	}
}
