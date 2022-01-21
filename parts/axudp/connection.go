package axudp

import (
	"github.com/rs/zerolog/log"
	"net"
	"sync"
	"time"
)

type Connection struct {
	Addr     *net.UDPAddr
	ttlTimer *time.Timer
}

var connectionTtl = time.Second * 10
var connectionMap = map[string]*Connection{}
var connectionLock = sync.RWMutex{}

func addOrUpdateConnection(pck *UDPServerPacket) error {
	connectionLock.RLock()
	addrString := pck.Addr.String()
	conn, ok := connectionMap[addrString]
	connectionLock.RUnlock()
	if !ok {
		connectionLock.Lock()
		conn = &Connection{
			Addr:     pck.Addr,
			ttlTimer: time.NewTimer(connectionTtl),
		}
		connectionMap[addrString] = conn
		go func() {
		MainLoop:
			for {
				select {
				case <-conn.ttlTimer.C:
					break MainLoop
				}
			}
			log.Debug().Str("addr", conn.Addr.String()).Msg("timeout")
			connectionLock.Lock()
			delete(connectionMap, addrString)
			connectionLock.Unlock()
		}()
		connectionLock.Unlock()
	}
	conn.ttlTimer.Reset(connectionTtl)
	return nil
}
