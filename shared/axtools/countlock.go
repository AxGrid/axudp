package axtools

import (
	"sync"
	"sync/atomic"
)

type CountLock struct {
	lock     sync.Mutex
	MaxCount int32
	count    int32
	locked   bool
}

type RWCountLock struct {
	lock     sync.RWMutex
	MaxCount int32
	count    int32
	locked   bool
}

func NewCountLock(maxCount int32) *CountLock {
	res := &CountLock{
		MaxCount: maxCount,
		lock:     sync.Mutex{},
		count:    0,
	}
	return res
}

func NewRWCountLock(maxCount int32) *RWCountLock {
	res := &RWCountLock{
		MaxCount: maxCount,
		lock:     sync.RWMutex{},
		count:    0,
	}
	return res
}

func (l *CountLock) IsLocked() bool {
	return l.locked
}
func (l *CountLock) Lock() {
	l.lock.Lock()
	l.locked = true
	cnt := atomic.AddInt32(&l.count, 1)
	if cnt < l.MaxCount {
		l.locked = false
		l.lock.Unlock()
	}
}
func (l *CountLock) Unlock() {
	atomic.AddInt32(&l.count, -1)
	if l.locked {
		l.locked = false
		l.lock.Unlock()
	}
}

func (l *RWCountLock) IsLocked() bool {
	return l.locked
}
func (l *RWCountLock) Lock() {
	l.lock.Lock()
	l.locked = true
	cnt := atomic.AddInt32(&l.count, 1)
	if cnt < l.MaxCount {
		l.locked = false
		l.lock.Unlock()
	}
}
func (l *RWCountLock) Unlock() {
	atomic.AddInt32(&l.count, -1)
	if l.locked {
		l.locked = false
		l.lock.Unlock()
	}
}
