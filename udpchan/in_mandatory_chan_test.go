package udpchan

import (
	pproto "github.com/axgrid/axudp/generated-sources/proto/axudp"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func TestInMandatoryChannelChunk(t *testing.T) {
	srvChan := make(chan []byte)
	sendChan := make(chan []byte)
	errorChan := make(chan error)

	go func() {
		for {
			select {
			case b := <-srvChan:
				log.Debug().Hex("payload", b).Msg("receive into Service")
			case b := <-sendChan:
				log.Debug().Hex("pck", b).Msg("Send to UDP (bytes)")
				var pck pproto.Packet
				assert.Nil(t, proto.Unmarshal(b, &pck))
				log.Debug().Hex("payload", b).Uint64("id", pck.Id).Msg("Send to UDP")

			case err := <-errorChan:
				log.Error().Err(err).Msg("error")
			}
		}
	}()
	packetWaitTTL = time.Millisecond * 10
	packetDoneTTL = time.Millisecond * 40

	imc := NewInMandatoryChannel(pproto.PacketMode_PM_MANDATORY, 100, srvChan, sendChan, errorChan, log.With().Logger())
	imc.Receive(&pproto.Packet{
		Id:         1,
		Parts:      booleansToBytes([]bool{true, false}),
		PartsCount: 2,
		Payload:    []byte{0, 1, 2},
	})

	imc.Receive(&pproto.Packet{
		Id:         1,
		Parts:      booleansToBytes([]bool{false, true}),
		PartsCount: 2,
		Payload:    []byte{4, 5},
	})
	time.Sleep(time.Millisecond * 50)
	assert.Equal(t, len(imc.incoming), 0)
}

func TestNewInMandatoryChannel(t *testing.T) {
	srvChan := make(chan []byte)
	sendChan := make(chan []byte)
	errorChan := make(chan error)

	go func() {
		for {
			select {
			case b := <-srvChan:
				log.Debug().Hex("payload", b).Msg("receive into Service")
			case b := <-sendChan:
				log.Debug().Hex("pck", b).Msg("Send to UDP (bytes)")
				var pck pproto.Packet
				assert.Nil(t, proto.Unmarshal(b, &pck))
				log.Debug().Hex("payload", b).Uint64("id", pck.Id).Msg("Send to UDP")

			case err := <-errorChan:
				log.Error().Err(err).Msg("error")
			}
		}
	}()

	packetWaitTTL = time.Millisecond * 10
	packetDoneTTL = time.Millisecond * 40
	imc := NewInMandatoryChannel(pproto.PacketMode_PM_MANDATORY, 100, srvChan, sendChan, errorChan, log.With().Logger())
	imc.Receive(&pproto.Packet{
		Id:         1,
		Parts:      booleansToBytes([]bool{true}),
		PartsCount: 1,
		Payload:    []byte{0, 1, 2},
	})
	time.Sleep(time.Millisecond * 50)
	assert.Equal(t, len(imc.incoming), 0)
}

func TestInOutMandatoryChannel(t *testing.T) {
	srvChan := make(chan []byte)
	sendChanI := make(chan []byte)
	sendChanO := make(chan []byte)
	errorChan := make(chan error)

	go func() {
		for {
			select {
			case b := <-srvChan:
				log.Debug().Hex("payload", b).Msg("receive into Service")
			case err := <-errorChan:
				log.Error().Err(err).Msg("error")
			}
		}
	}()

	packetWaitTTL = time.Millisecond * 10
	packetDoneTTL = time.Millisecond * 40

	imc := NewInMandatoryChannel(pproto.PacketMode_PM_MANDATORY, 100, srvChan, sendChanI, errorChan, log.With().Logger())
	omc := NewOutMandatoryChannel(100, sendChanO, errorChan, log.With().Logger())

	go func() {
		for {
			select {
			case send := <-sendChanI:
				var pck pproto.Packet
				assert.Nil(t, proto.Unmarshal(send, &pck))
				log.Debug().Uint64("id", pck.Id).Str("tag", pck.Tag.String()).Msg("chan-i recv")
				omc.Response(&pck)
			case send := <-sendChanO:
				var pck pproto.Packet
				assert.Nil(t, proto.Unmarshal(send, &pck))
				log.Debug().Uint64("id", pck.Id).Str("tag", pck.Tag.String()).Bools("parts", bytesToBooleans(pck.Parts)[:pck.PartsCount]).Uint32("partCount", pck.PartsCount).Hex("payload", pck.Payload).Msg("chan-o recv")
				imc.Receive(&pck)
			}
		}
	}()

	omc.Send([]byte{1, 2, 3, 4, 5})
	time.Sleep(time.Millisecond * 10)

	assert.Equal(t, len(imc.incoming), 1)
	assert.Equal(t, len(omc.outgoing), 0)
	time.Sleep(time.Millisecond * 50)
	assert.Equal(t, len(imc.incoming), 0)
	assert.Equal(t, len(omc.outgoing), 0)

}
