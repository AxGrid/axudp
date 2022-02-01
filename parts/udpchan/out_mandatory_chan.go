package udpchan

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"sync"
	"sync/atomic"
)

type IResponse interface {
	Response(pck *pproto.Packet)
}

type OutMandatoryChannel struct {
	id           uint64
	lock         sync.RWMutex
	mtu          int
	mode         pproto.PacketMode
	sendChan     chan []byte
	responseChan chan HolderResponse
	errorChan    chan error
	closeChan    chan bool
	outgoing     map[uint64]*OutHolder
}

func NewOutMandatoryChannel(mtu int, sendChan chan []byte, errorChan chan error) *OutMandatoryChannel {
	res := &OutMandatoryChannel{
		lock:         sync.RWMutex{},
		mtu:          mtu,
		sendChan:     sendChan,
		errorChan:    errorChan,
		mode:         pproto.PacketMode_PM_MANDATORY,
		closeChan:    make(chan bool, 1),
		responseChan: make(chan HolderResponse),
		outgoing:     map[uint64]*OutHolder{},
	}

	go func() {
		for {
			select {
			case <-res.closeChan:
				for _, h := range res.outgoing {
					h.Stop()
				}
				return
			case r := <-res.responseChan:
				res.lock.Lock()
				if r.err != nil {
					for _, h := range res.outgoing {
						h.Stop()
					}
					res.errorChan <- r.err
					delete(res.outgoing, r.id)
					return
				} else {
					h, ok := res.outgoing[r.id]
					if ok {
						h.Stop()
						delete(res.outgoing, r.id)
					}
				}
				res.lock.Unlock()
			}
		}
	}()
	return res
}

func (c *OutMandatoryChannel) Stop() {
	c.closeChan <- true
}

func (c *OutMandatoryChannel) Send(payload []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()
	id := atomic.AddUint64(&c.id, 1)
	var err error
	c.outgoing[id], err = NewOutHolder(id, payload, c.mtu, c.responseChan, c.mode, c.sendChan)
	if err != nil {
		c.errorChan <- err
		return
	}
	c.outgoing[id].Start()
	c.outgoing[id].send()
}

func (c *OutMandatoryChannel) Response(pck *pproto.Packet) {
	c.lock.RLock()
	h, ok := c.outgoing[pck.Id]
	c.lock.RUnlock()
	if !ok {
		return
	}
	h.Response(pck)
}

type OutMandatoryConsistentChannel struct {
	id           uint64
	lock         sync.RWMutex
	mtu          int
	mode         pproto.PacketMode
	sendChan     chan []byte
	responseChan chan HolderResponse
	errorChan    chan error
	closeChan    chan bool
	outgoing     *OutHolder
	queue        [][]byte
}

func NewOutMandatoryConsistentChannel(mtu int, sendChan chan []byte, errorChan chan error) *OutMandatoryConsistentChannel {
	res := &OutMandatoryConsistentChannel{
		lock:         sync.RWMutex{},
		mtu:          mtu,
		sendChan:     sendChan,
		errorChan:    errorChan,
		mode:         pproto.PacketMode_PM_MANDATORY,
		closeChan:    make(chan bool, 1),
		responseChan: make(chan HolderResponse),
		outgoing:     nil,
		queue:        [][]byte{},
	}

	go func() {
		for {
			select {
			case <-res.closeChan:
				if res.outgoing != nil {
					res.outgoing.Stop()
				}
				return
			case r := <-res.responseChan:
				res.lock.Lock()
				if r.err != nil {
					if res.outgoing != nil {
						res.outgoing.Stop()
					}
					res.errorChan <- r.err
					return
				} else {
					if res.outgoing != nil && res.outgoing.id == r.id {
						res.outgoing.Stop()
						res.outgoing = nil
						go func() { res.next() }()
					}
				}
				res.lock.Unlock()
			}
		}
	}()
	return res
}

func (c *OutMandatoryConsistentChannel) next() {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.outgoing != nil {
		return
	}
	if len(c.queue) == 0 {
		return
	}
	var err error
	payload := c.queue[0]
	c.queue = c.queue[1:]
	id := atomic.AddUint64(&c.id, 1)
	c.outgoing, err = NewOutHolder(id, payload, c.mtu, c.responseChan, c.mode, c.sendChan)
	if err != nil {
		c.outgoing = nil
		c.errorChan <- err
		return
	}
	c.outgoing.Start()
	c.outgoing.send()
}

func (c *OutMandatoryConsistentChannel) Send(payload []byte) {
	c.lock.Lock()
	c.queue = append(c.queue, payload)
	c.lock.Unlock()
	c.next()
}

func (c *OutMandatoryConsistentChannel) Response(pck *pproto.Packet) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.outgoing == nil {
		return
	}
	c.outgoing.Response(pck)
}

func (c *OutMandatoryConsistentChannel) Stop() {
	c.closeChan <- true
}
