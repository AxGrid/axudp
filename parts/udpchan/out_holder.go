package udpchan

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

var packetRetryTTL = time.Millisecond * 1200
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
}

func NewOutHolder(id uint64, payload []byte, mtu int, responseChan chan HolderResponse, mode pproto.PacketMode, sendChan chan []byte) (*OutHolder, error) {
	chunks := getChunk(payload, mtu)
	partsCount := len(chunks)
	outPacket := make([][]byte, partsCount)
	log.Trace().Str("class", "out_holder").Int("parts", partsCount).Int("mtu", mtu).Hex("hex", payload).Msg("create out holder")
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

	return &OutHolder{
		id:           id,
		responseChan: responseChan,
		sendChan:     sendChan,
		parts:        make([]bool, partsCount),
		outs:         outPacket,
		closeChan:    nil,
	}, nil
}

func (h *OutHolder) Response(pck *pproto.Packet) {
	h.lock.Lock()
	mergeBoolLists(&h.parts, bytesToBooleans(pck.Parts))
	h.lock.Unlock()
	if allBooleans(h.parts) {
		h.done = true
		h.responseChan <- HolderResponse{id: h.id}
		h.Stop()
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
					h.send()
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
