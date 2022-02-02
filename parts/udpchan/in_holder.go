package udpchan

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"fmt"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

var packetWaitTTL = time.Millisecond * 1000
var packetDoneTTL = time.Millisecond * 5000

type InHolder struct {
	id           uint64
	parts        []bool
	partsCount   int
	payload      [][]byte
	done         bool
	mode         pproto.PacketMode
	sendChan     chan []byte
	responseChan chan HolderResponse // Канал для сообщения об ошибках и завершении отправки
	closeChan    chan bool
	lock         sync.Mutex
	timer        *time.Timer
	log          zerolog.Logger
}

func NewInHolder(id uint64, partsCount int, responseChan chan HolderResponse, mode pproto.PacketMode, sendChan chan []byte, llog zerolog.Logger) *InHolder {
	return &InHolder{
		id:           id,
		responseChan: responseChan,
		parts:        make([]bool, partsCount),
		partsCount:   partsCount,
		payload:      make([][]byte, partsCount),
		mode:         mode,
		sendChan:     sendChan,
		closeChan:    nil,
		log:          llog.With().Str("mode", mode.String()).Uint64("id", id).Str("holder", "in_holder").Int("partCount", partsCount).Logger(),
	}
}

func (h *InHolder) Receive(pck *pproto.Packet) {

	index := booleanIndex(bytesToBooleans(pck.Parts))
	if index < 0 || index >= h.partsCount {
		h.responseChan <- HolderResponse{id: h.id, err: fmt.Errorf("wrong part index %d >= %d", index, h.partsCount)}
		return
	}

	h.lock.Lock()
	defer h.lock.Unlock()

	if h.parts[index] {
		h.response()
		return
	}

	if h.done {
		return
	}
	if h.timer != nil {
		h.timer.Reset(packetWaitTTL)
	}

	h.parts[index] = true
	h.payload[index] = pck.Payload

	h.response()
	if allBooleans(h.parts) {
		if h.timer != nil {
			h.timer.Reset(packetDoneTTL)
		}
		h.log.Trace().Msg("packet done")
		h.done = true
		h.responseChan <- HolderResponse{id: h.id, payload: bytesJoin(h.payload...)}
	}
}

func (h *InHolder) Stop() {
	if h.closeChan != nil {
		h.log.Trace().Msg("holder stop")
		h.closeChan <- true
	}
}

func (h *InHolder) Start() {
	if h.closeChan != nil {
		return
	}
	h.closeChan = make(chan bool, 1)
	h.timer = time.NewTimer(packetWaitTTL)
	go func() {
		for {
			select {
			case <-h.closeChan:
				return
			case <-h.timer.C:
				if h.done {
					h.responseChan <- HolderResponse{id: h.id}
				} else {
					h.responseChan <- HolderResponse{id: h.id, err: fmt.Errorf("timeout in packet %d", h.id)}
				}
				return
			}
		}
	}()
}

func (h *InHolder) response() {
	p := &pproto.Packet{
		Id:         h.id,
		Parts:      booleansToBytes(h.parts),
		PartsCount: uint32(h.partsCount),
		Mode:       h.mode,
		Tag:        pproto.PacketTag_PT_DONE,
	}
	bytes, err := proto.Marshal(p)
	if err != nil {
		h.responseChan <- HolderResponse{id: h.id, err: err}
		return
	}
	go func() { h.sendChan <- bytes }()
}
