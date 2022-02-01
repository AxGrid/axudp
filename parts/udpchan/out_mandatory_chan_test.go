package udpchan

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewOutMandatoryChannel(t *testing.T) {
	sendChan := make(chan []byte)
	errorChan := make(chan error)
	var resendCount uint64 = 0
	var er error
	go func() {
		for {
			select {
			case send := <-sendChan:
				atomic.AddUint64(&resendCount, 1)
				log.Debug().Hex("send", send).Msg("Send to UDP")
			case err := <-errorChan:
				er = err
				log.Error().Err(err).Msg("fail of holder")
			}
		}
	}()
	packetRetryTTL = time.Millisecond * 20
	packetRetryCount = 4
	omc := NewOutMandatoryChannel(50, sendChan, errorChan)
	omc.Send([]byte{0, 1, 2, 3, 4, 5})
	time.Sleep(time.Millisecond * 50)
	assert.Nil(t, er)
	assert.Equal(t, resendCount, uint64(3))
	time.Sleep(time.Millisecond * 50)
	assert.Equal(t, resendCount, uint64(4))
	assert.NotNil(t, er)
	assert.Equal(t, len(omc.outgoing), 0)
}

func TestNewOutMandatoryChunkChannel(t *testing.T) {
	sendChan := make(chan []byte)
	errorChan := make(chan error)
	var resendCount uint64 = 0
	var er error
	go func() {
		for {
			select {
			case send := <-sendChan:
				atomic.AddUint64(&resendCount, 1)
				log.Debug().Hex("send", send).Msg("Send to UDP")
			case err := <-errorChan:
				er = err
				log.Error().Err(err).Msg("fail of holder")
			}
		}
	}()
	packetRetryTTL = time.Millisecond * 20
	packetRetryCount = 6
	omc := NewOutMandatoryChannel(3, sendChan, errorChan)
	omc.Send([]byte{0, 1, 2, 3, 4})
	time.Sleep(time.Millisecond * 25)
	assert.Nil(t, er)
	assert.Equal(t, resendCount, uint64(4))
	omc.Response(&pproto.Packet{
		Id:    1,
		Parts: booleansToBytes([]bool{true, false}),
	})
	time.Sleep(time.Millisecond * 30)
	omc.Response(&pproto.Packet{
		Id:    1,
		Parts: booleansToBytes([]bool{true, true}),
	})
	time.Sleep(time.Millisecond * 50)
	assert.Equal(t, len(omc.outgoing), 0)
}

func TestNewOutMandatoryConsistentChannel(t *testing.T) {
	sendChan := make(chan []byte)
	errorChan := make(chan error)
	var resendCount uint64 = 0
	var er error
	go func() {
		for {
			select {
			case send := <-sendChan:
				atomic.AddUint64(&resendCount, 1)
				var pck pproto.Packet
				assert.Nil(t, proto.Unmarshal(send, &pck))

				log.Debug().Hex("send", pck.Payload).Uint64("id", pck.Id).Bools("parts", bytesToBooleans(pck.Parts)[:pck.PartsCount]).Msg("Send to UDP")
			case err := <-errorChan:
				er = err
				log.Error().Err(err).Msg("fail of holder")
			}
		}
	}()
	packetRetryTTL = time.Millisecond * 10
	packetRetryCount = 4
	omc := NewOutMandatoryConsistentChannel(50, sendChan, errorChan)
	omc.Send([]byte{0, 1, 2, 3, 4, 5})
	omc.Send([]byte{6, 7, 8})
	time.Sleep(time.Millisecond * 25)
	omc.Response(&pproto.Packet{
		Id:    1,
		Parts: booleansToBytes([]bool{true}),
	})
	assert.Equal(t, resendCount, uint64(3))
	time.Sleep(time.Millisecond * 36)
	omc.Response(&pproto.Packet{
		Id:    2,
		Parts: booleansToBytes([]bool{true}),
	})
	assert.Equal(t, resendCount, uint64(4+3))
	assert.Nil(t, er)
}
