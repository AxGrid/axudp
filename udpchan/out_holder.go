package udpchan

import (
	"errors"
	pproto "github.com/axgrid/axudp/generated-sources/proto/axudp"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

var packetRetryTTL = time.Millisecond * 100
var packetRetryCount = 5

var timeOutRetryMap = map[int]time.Duration{
	0:  time.Millisecond * 100,
	5:  time.Millisecond * 200,
	10: time.Millisecond * 500,
	20: time.Millisecond * 1000,
}

func getTimeout(retryCount int) time.Duration {
	var res time.Duration = 0
	for k, v := range timeOutRetryMap {
		if retryCount >= k {
			res = v
		}
	}
	return res
}

type HolderResponse struct {
	err     error
	id      uint64
	payload []byte
}

type OutHolder struct {
	id           uint64
	parts        []bool
	outs         [][]byte
	retryCount   int
	done         bool
	responseChan chan HolderResponse // Канал для сообщения об ошибках и завершении отправки
	sendChan     chan []byte
	closeChan    chan bool
	lock         sync.Mutex
	log          zerolog.Logger
	tm           *time.Ticker
}

func NewOutHolder(id uint64, payload []byte, mtu int, responseChan chan HolderResponse, mode pproto.PacketMode, sendChan chan []byte, llog zerolog.Logger) (*OutHolder, error) {
	chunks := getChunk(payload, mtu)
	partsCount := len(chunks)
	outPacket := make([][]byte, partsCount)
	tag := pproto.PacketTag_PT_PAYLOAD
	var err error
	for i, chunk := range chunks {
		parts := make([]bool, partsCount)
		parts[i] = true
		rpck := &pproto.Packet{
			Id:         id,
			Parts:      booleansToBytes(parts),
			PartsCount: uint32(partsCount),
			Payload:    chunk,
			Mode:       mode,
			Mtu:        uint32(mtu),
			Tag:        tag,
		}
		outPacket[i], err = proto.Marshal(rpck)
		if err != nil {
			return nil, err
		}
	}

	res := &OutHolder{
		id:           id,
		responseChan: responseChan,
		sendChan:     sendChan,
		parts:        make([]bool, partsCount),
		outs:         outPacket,
		closeChan:    nil,
		log:          llog.With().Str("mode", mode.String()).Int("size", len(payload)).Int("mtu", mtu).Uint64("id", id).Str("holder", "out_holder").Int("partCount", partsCount).Logger(),
	}
	res.log.Trace().Msg("create holder")
	return res, nil
}

func (h *OutHolder) Response(pck *pproto.Packet) {
	if pck.Id != h.id {
		return
	}
	h.lock.Lock()
	defer h.lock.Unlock()
	h.retryCount = 0
	h.tm.Reset(timeOutRetryMap[0])
	mergeBoolLists(&h.parts, bytesToBooleans(pck.Parts))

	if allBooleans(h.parts) && !h.done {
		h.log.Trace().Bools("parts", h.parts).Msg("holder done")
		h.done = true
		h.responseChan <- HolderResponse{id: h.id}
		h.Stop()
	} else {
		h.log.Trace().Bools("parts", h.parts).Msg("holder collect")
	}
}

func (h *OutHolder) Stop() {
	if h.closeChan != nil {
		h.closeChan <- true
	}
}

func (h *OutHolder) Start() {
	if h.closeChan != nil {
		return
	}
	h.closeChan = make(chan bool, 1)
	h.tm = time.NewTicker(packetRetryTTL)
	go func() {
	MainLoop:
		for {
			select {
			case <-h.closeChan:
				break MainLoop
			case <-h.tm.C:
				h.retryCount++
				ttl := getTimeout(h.retryCount)
				if ttl == 0 {
					h.responseChan <- HolderResponse{
						id:  h.id,
						err: errors.New("max retry count"),
					}
					break MainLoop
				} else {
					h.tm.Reset(ttl)
					h.log.Trace().Bools("parts", h.parts).Msg("retry")
					go func() { h.send() }()
				}
			}
		}
		h.tm.Stop()
	}()
}

func (h *OutHolder) send() {
	for i, part := range h.parts {
		if part {
			continue
		}
		h.sendChan <- h.outs[i]
	}
}
