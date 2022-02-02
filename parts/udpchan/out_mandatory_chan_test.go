package udpchan

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"os"
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
	omc := NewOutMandatoryChannel(50, sendChan, errorChan, log.With().Logger())
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
	omc := NewOutMandatoryChannel(3, sendChan, errorChan, log.With().Logger())
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
	omc := NewOutMandatoryConsistentChannel(50, sendChan, errorChan, log.With().Logger())
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

func TestNewOutMandatoryConsistentMultiChannel(t *testing.T) {

	sendChan := make(chan []byte)
	errorChan := make(chan error)
	var resendCount uint64 = 0
	var er error
	omc := NewOutMandatoryConsistentChannel(10, sendChan, errorChan, log.With().Logger())
	count := 0
	go func() {
		var id uint64 = 0
		var parts []bool
		for {
			select {
			case send := <-sendChan:
				atomic.AddUint64(&resendCount, 1)
				var pck pproto.Packet
				assert.Nil(t, proto.Unmarshal(send, &pck))
				log.Debug().Hex("send", pck.Payload).Uint64("id", pck.Id).Bools("parts", bytesToBooleans(pck.Parts)[:pck.PartsCount]).Msg("Send to UDP")
				//omc.Response(pck)
				if id != pck.Id {
					count++
					parts = make([]bool, pck.PartsCount)
					id = pck.Id
				}
				ind := booleanIndex(bytesToBooleans(pck.Parts))
				parts[ind] = true
				log.Debug().Uint64("id", pck.Id).Bools("parts", parts).Msg("Response")
				omc.Response(&pproto.Packet{
					Id:         id,
					Mode:       pproto.PacketMode_PM_MANDATORY_CONSISTENTLY,
					Tag:        pproto.PacketTag_PT_DONE,
					Parts:      booleansToBytes(parts),
					PartsCount: pck.PartsCount,
				})
			case err := <-errorChan:
				er = err
				log.Error().Err(err).Msg("fail of holder")
			}
		}
	}()

	packetRetryTTL = time.Millisecond * 100
	packetRetryCount = 4

	for i := 1; i < 50; i++ {
		omc.Send(make([]byte, i+12))
	}
	time.Sleep(time.Millisecond * 2000)
	assert.Nil(t, er)
	assert.Equal(t, len(omc.queue), 0)
	assert.Nil(t, omc.outgoing)
	assert.Equal(t, count, 49)
}

func TestNewOutToInMandatoryConsistentMultiChannel(t *testing.T) {
	sendToInChan := make(chan []byte)
	sendResponseChan := make(chan []byte)
	errorChan := make(chan error)
	serviceChan := make(chan []byte)
	omc := NewOutMandatoryConsistentChannel(10, sendToInChan, errorChan, log.With().Logger())
	imc := NewInMandatoryChannel(pproto.PacketMode_PM_MANDATORY_CONSISTENTLY, 10, serviceChan, sendResponseChan, errorChan, log.With().Logger())

	var resp []int
	var errCount int
	go func() {
		for {
			select {
			case send := <-sendToInChan:
				var pck pproto.Packet
				proto.Unmarshal(send, &pck)
				imc.Receive(&pck)
			case send := <-sendResponseChan:
				var pck pproto.Packet
				proto.Unmarshal(send, &pck)
				omc.Response(&pck)
			case serv := <-serviceChan:
				resp = append(resp, len(serv))
			case err := <-errorChan:
				log.Error().Err(err).Msg("error in channels")
				errCount++
			}
		}
	}()

	packetRetryTTL = time.Millisecond * 100
	packetRetryCount = 4

	for i := 0; i < 50; i++ {
		omc.Send(make([]byte, i+12))
	}

	time.Sleep(time.Millisecond * 2000)
	assert.Equal(t, errCount, 0)
	assert.Equal(t, len(resp), 50)
	assert.Equal(t, len(omc.queue), 0)
	assert.Nil(t, omc.outgoing)
}

func TestMain(m *testing.M) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05,000"}).Level(zerolog.TraceLevel)
	code := m.Run()
	os.Exit(code)
}
