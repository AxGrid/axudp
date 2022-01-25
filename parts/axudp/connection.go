package axudp

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"github.com/rs/zerolog/log"
	"net"
	"sync"
	"time"
)

type Connection struct {
	Addr      *net.UDPAddr
	ttlTimer  *time.Timer
	InStream  *InStream
	OutStream *OutStream
	InChannel chan []byte
}

var connectionTtl = time.Second * 10
var connectionMap = map[string]*Connection{}
var connectionLock = sync.RWMutex{}

func newConnection(addr *net.UDPAddr, outChan chan *pproto.Packet) *Connection {
	servChan := make(chan []byte)
	return &Connection{
		Addr:      addr,
		ttlTimer:  time.NewTimer(connectionTtl),
		InStream:  NewInStream(outChan, servChan),
		OutStream: NewOutStream(outChan),
	}
}

func (c *Connection) receive(pck *pproto.Packet) {
	c.ttlTimer.Reset(connectionTtl)
	c.InStream.receive(pck)
}

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
