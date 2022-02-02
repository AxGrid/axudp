package udpchan

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"github.com/rs/zerolog"
	"sync"
)

type InMandatoryChannel struct {
	lock         sync.RWMutex
	mode         pproto.PacketMode
	mtu          int
	sendChan     chan []byte
	errorChan    chan error
	closeChan    chan bool
	serviceChan  chan []byte
	responseChan chan HolderResponse
	incoming     map[uint64]*InHolder
	log          zerolog.Logger
}

func NewInMandatoryChannel(mode pproto.PacketMode, mtu int, serviceChan chan []byte, sendChan chan []byte, errorChan chan error, llog zerolog.Logger) *InMandatoryChannel {
	res := &InMandatoryChannel{
		lock:         sync.RWMutex{},
		mtu:          mtu,
		mode:         mode,
		sendChan:     sendChan,
		errorChan:    errorChan,
		closeChan:    make(chan bool, 1),
		serviceChan:  serviceChan,
		responseChan: make(chan HolderResponse),
		incoming:     map[uint64]*InHolder{},
		log:          llog.With().Str("mode", mode.String()).Str("chan", "in_mandatory_channel").Logger(),
	}

	go func() {
		for {
			select {
			case <-res.closeChan:
				for _, h := range res.incoming {
					h.Stop()
				}
				return
			case r := <-res.responseChan:
				switch {
				case r.err != nil:
					res.errorChan <- r.err
					return
				case r.payload != nil:
					serviceChan <- r.payload
					break
				default:
					res.delete(r.id)
					return
				}

			}
		}
	}()
	return res
}

func (c *InMandatoryChannel) Receive(pck *pproto.Packet) {
	c.lock.Lock()
	defer c.lock.Unlock()
	holder, ok := c.incoming[pck.Id]
	if !ok {

		holder = NewInHolder(pck.Id, int(pck.PartsCount), c.responseChan, c.mode, c.sendChan, c.log)
		c.incoming[pck.Id] = holder
		holder.Start()
	}
	holder.Receive(pck)
}

func (c *InMandatoryChannel) delete(id uint64) {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.incoming[id]
	if !ok {
		return
	}
	delete(c.incoming, id)
}

func (c *InMandatoryChannel) Stop() {
	if c.closeChan != nil {
		c.closeChan <- true
	}
}
