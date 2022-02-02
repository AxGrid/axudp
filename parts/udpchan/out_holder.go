package udpchan

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"errors"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

var packetRetryTTL = time.Millisecond * 100
var packetRetryCount = 5

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
	go func() {
		tm := time.NewTicker(packetRetryTTL)
	MainLoop:
		for {
			select {
			case <-h.closeChan:
				break MainLoop
			case <-tm.C:
				h.retryCount++
				if h.retryCount >= packetRetryCount {
					h.responseChan <- HolderResponse{
						id:  h.id,
						err: errors.New("max retry count"),
					}
					break MainLoop
				} else {
					h.log.Trace().Bools("parts", h.parts).Msg("retry")
					go func() { h.send() }()
				}
			}
		}
		tm.Stop()
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
