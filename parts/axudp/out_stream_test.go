package axudp

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"time"

	"testing"
)

func TestOutOptional(t *testing.T) {
	S := make(chan *pproto.Packet)

	go func() {
		for i := 0; i < 2; i++ {
			pck := <-S
			log.Debug().Hex("payload", pck.Payload).Uint64("id", pck.Id).Uint32("part", pck.PartsCount).Msg("receive")
			switch i {
			case 0:
				assert.Equal(t, pck.Payload, []byte{1, 2, 3})
			case 1:
				assert.Equal(t, pck.Payload, []byte{4, 5})
			}
		}
	}()
	c := newOutOptionalStream(S, pproto.PacketMode_PM_OPTIONAL)
	c.mtu = 3
	c.hold([]byte{1, 2, 3, 4, 5})
	time.Sleep(time.Millisecond * 10)
}

func TestOutMandatory(t *testing.T) {
	S := make(chan *pproto.Packet)
	c := newOutMandatoryStream(S)
	go func() {
		i := 0
		for {
			pck := <-S
			log.Debug().Hex("payload", pck.Payload).Uint64("id", pck.Id).Uint32("part", pck.PartsCount).Msg("receive")
			i++
			if i > 3 {
				c.response(&pproto.Packet{
					Id:         pck.Id,
					Parts:      booleansToBytes([]bool{true, true}),
					PartsCount: pck.PartsCount,
				})
			}
		}
	}()
	retryTTL = time.Millisecond * 10
	c.mtu = 3
	c.hold([]byte{1, 2, 3, 4, 5})
	c.hold([]byte{8, 9})
	time.Sleep(time.Millisecond * 110)
}

func TestOutMandatoryConsistently(t *testing.T) {
	S := make(chan *pproto.Packet)
	c := newOutMandatoryConsistentlyStream(S)
	go func() {
		i := 0
		for {
			pck := <-S
			log.Debug().Hex("payload", pck.Payload).Uint64("id", pck.Id).Uint32("part", pck.PartsCount).Msg("receive")
			i++
			if i > 3 {
				c.response(&pproto.Packet{
					Id:         pck.Id,
					Parts:      booleansToBytes([]bool{true, true}),
					PartsCount: pck.PartsCount,
				})
			}
		}
	}()

	retryTTL = time.Millisecond * 10
	c.mtu = 3
	c.hold([]byte{1, 2, 3, 4, 5})
	c.hold([]byte{8, 9})
	time.Sleep(time.Millisecond * 110)
}
