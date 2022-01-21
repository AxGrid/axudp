package axtools

import (
	"github.com/rs/zerolog/log"
	"sync"
)

var idLocks = map[uint64]*sync.Mutex{}
var mapLock = sync.RWMutex{}

func LockID(id uint64) {
	if id == 0 {
		return
	}
	mapLock.RLock()
	lock, ok := idLocks[id]
	mapLock.RUnlock()
	if !ok {
		lock = &sync.Mutex{}
		mapLock.Lock()
		idLocks[id] = lock
		mapLock.Unlock()
	}
	lock.Lock()
	log.Trace().Uint64("uid", id).Msg("lock")
}

func UnlockID(id uint64) {
	if id == 0 {
		return
	}
	//mapLock.Lock()
	//defer mapLock.Unlock()

	mapLock.RLock()
	lock, ok := idLocks[id]
	mapLock.RUnlock()
	if ok {
		lock.Unlock()
		log.Trace().Uint64("uid", id).Msg("unlock")
		mapLock.Lock()
		delete(idLocks, id)
		mapLock.Unlock()
	}
}
