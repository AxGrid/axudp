package axudp

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"errors"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

var waitPacketTTL = time.Millisecond * 1000
var sendPacketTTL = time.Millisecond * 5000
var packetRetryTTL = time.Millisecond * 100

type iChannel interface {
	delete(id uint64)
	done(id uint64)
}

type incomingChannel struct {
	holders map[uint64]*incomingHolder
	lastId  uint64
	maxId   uint64
	mode    pproto.PacketMode
	lock    sync.Mutex
	C       chan []byte         // Каннал для отправки результата в сервис (payload)
	D       chan *pproto.Packet //Канал для отправки ответа
}

func newIncomingChannel(mode pproto.PacketMode, responseChannel chan *pproto.Packet) *incomingChannel {
	return &incomingChannel{
		holders: map[uint64]*incomingHolder{},
		lastId:  0,
		mode:    mode,
		lock:    sync.Mutex{},
		C:       make(chan []byte),
		D:       responseChannel,
	}
}

func (ic *incomingChannel) hold(pck *pproto.Packet) error {
	ic.lock.Lock()
	defer ic.lock.Unlock()
	if ic.maxId < pck.Id {
		ic.maxId = pck.Id
	}
	switch ic.mode {
	case pproto.PacketMode_PM_OPTIONAL, pproto.PacketMode_PM_MANDATORY:
		ih, ok := ic.holders[pck.Id]
		if !ok {
			ih = newIncomingHolder(pck, ic)
			ic.holders[pck.Id] = ih
		}
		if err := ih.hold(pck); err != nil {
			return err
		}
		ic.D <- ih.status(pproto.PacketTag_PT_DONE)
		if ih.done {
			go ic.done(ih.id)
		}
	case pproto.PacketMode_PM_OPTIONAL_CONSISTENTLY:
		if pck.Id < ic.lastId {
			ic.D <- &pproto.Packet{Id: pck.Id, Tag: pproto.PacketTag_PT_OUTDATE}
			return nil
		}
		ih := newIncomingHolder(pck, ic)
		if pck.Id > ic.lastId {
			if ic.lastId > 0 {
				go ic.delete(ic.lastId)
			}
			ic.lastId = pck.Id
			ic.holders = map[uint64]*incomingHolder{pck.Id: ih}
		}
		if err := ih.hold(pck); err != nil {
			return err
		}
		ic.D <- ih.status(pproto.PacketTag_PT_DONE)
		if ih.done {
			go ic.done(ih.id)
		}
	case pproto.PacketMode_PM_MANDATORY_CONSISTENTLY:
		ih, ok := ic.holders[pck.Id]
		if !ok {
			ih = newIncomingHolder(pck, ic)
			ic.holders[pck.Id] = ih
		}
		if err := ih.hold(pck); err != nil {
			return err
		}
		ic.D <- ih.status(pproto.PacketTag_PT_DONE)
		if ih.done {
			log.Debug().Uint64("id", ih.id).Msg("IS DONE")
			for i := ic.lastId + 1; i <= ic.maxId; i++ {
				ih, ok = ic.holders[i]
				if !ok || !ih.done {
					log.Debug().Uint64("id", i).Msg("not found! break")
					break
				}
				log.Debug().Uint64("id", i).Msg("found!")
				ic.lastId = i
				ic.doneUnsafe(i)
			}
		}
	}
	return nil
}

func (ic *incomingChannel) done(id uint64) {
	ic.lock.Lock()
	defer ic.lock.Unlock()
	ic.doneUnsafe(id)
}

func (ic *incomingChannel) doneUnsafe(id uint64) {
	holder, ok := ic.holders[id]
	if !ok {
		return
	}
	log.Trace().Uint64("id", id).Msg("done holder")
	ic.C <- bytesJoin(holder.payloads...)
	holder.exit <- true
}

func (ic *incomingChannel) delete(id uint64) {
	ic.lock.Lock()
	defer ic.lock.Unlock()
	log.Trace().Uint64("id", id).Msg("delete holder")
	delete(ic.holders, id)
}

type incomingHolder struct {
	id         uint64
	parts      []bool
	partsCount int
	done       bool
	payloads   [][]byte
	tag        pproto.PacketTag
	exit       chan bool
	//parent     iChannel
}

func newIncomingHolder(pck *pproto.Packet, parent iChannel) *incomingHolder {
	res := &incomingHolder{
		id:         pck.Id,
		partsCount: int(pck.PartsCount),
		parts:      make([]bool, pck.PartsCount),
		payloads:   make([][]byte, pck.PartsCount),
		tag:        pck.Tag,
		done:       false,
		exit:       make(chan bool, 1),
		//parent:     parent,
	}
	go func() {
		tm := time.NewTimer(waitPacketTTL)
	MainLoop:
		for {
			select {
			case <-res.exit:
				break MainLoop
			case <-tm.C:
				log.Trace().Uint64("id", res.id).Msg("timeout holder")
				break MainLoop
			}
		}
		parent.delete(res.id)
	}()
	return res
}

func (ih *incomingHolder) hold(pck *pproto.Packet) error {
	pckBoolMask := bytesToBooleans(pck.Parts)
	mergeBoolLists(&ih.parts, pckBoolMask)
	index := booleanIndex(pckBoolMask)
	if index < 0 || index >= ih.partsCount {
		return errors.New("wrong part index")
	}
	ih.payloads[index] = pck.Payload
	ih.done = allBooleans(ih.parts)
	return nil
}

func (ih *incomingHolder) status(tag pproto.PacketTag) *pproto.Packet {
	return &pproto.Packet{
		Id:         ih.id,
		Tag:        tag,
		Parts:      booleansToBytes(ih.parts),
		PartsCount: uint32(ih.partsCount),
	}
}

type outgoingChannel struct {
	holders map[uint64]*incomingHolder
	lastId  uint64
	id      uint64
	mode    pproto.PacketMode
	lock    sync.Mutex
	R       chan *pproto.Packet // канал для получения ответа о доставке incomingChannel.D
}

type outgoingHolder struct {
	id         uint64
	parts      []bool
	partsCount int
	done       bool
	payloads   [][]byte
	tag        pproto.PacketTag
	exit       chan bool
	//parent     iChannel
	ttlTm *time.Timer
}

func newOutgoingHolder(parent iChannel, payload []byte, mtu int) *outgoingHolder {
	chank := getChunk(payload, mtu)
	parts := len(chank)
	res := &outgoingHolder{
		parts:      make([]bool, parts),
		partsCount: parts,
		done:       false,
		payloads:   chank,
		tag:        pproto.PacketTag_PT_PAYLOAD,
		exit:       make(chan bool, 1),
		//parent:     parent,
		ttlTm: time.NewTimer(sendPacketTTL),
	}

	go func() {
		retryTm := time.NewTicker(packetRetryTTL)

	MainLoop:
		for {
			select {
			case <-retryTm.C:
				//:TODO Retry
			case <-res.exit:
				break MainLoop
			case <-res.ttlTm.C:
				break MainLoop
			}
		}
		parent.delete(res.id)
	}()
	return res
}
