package axtools

import "sync"

type PubSub struct {
	listeners map[*Listener]bool
	lock      sync.RWMutex
	closed    bool
}

type Listener func([]byte)

func NewPubSub() *PubSub {
	return &PubSub{
		listeners: map[*Listener]bool{},
		lock:      sync.RWMutex{},
		closed:    false,
	}
}

func (c *PubSub) Close() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.closed = true
	c.listeners = nil
}

func (c *PubSub) Sub(f *Listener) {
	if c.closed {
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	c.listeners[f] = true
}

func (c *PubSub) UnSub(f *Listener) {
	if c.closed {
		return
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.listeners, f)
}

func (c *PubSub) Pub(data []byte) {
	if c.closed {
		return
	}
	c.lock.RLock()
	defer c.lock.RUnlock()
	for l := range c.listeners {
		if l != nil {
			(*l)(data)
		}
	}
}
