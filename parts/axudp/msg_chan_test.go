package axudp

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOptionalChannel(t *testing.T) {

	resp := make(chan *pproto.Packet)

	c := newIncomingChannel(pproto.PacketMode_PM_OPTIONAL, resp)

	go func() {
		for {
			r := <-resp
			log.Debug().Interface("pck", r).Msg("response")
		}
	}()

	go func() {
		data := <-c.C
		assert.EqualValues(t, data, []byte{0, 1, 2, 3, 4, 5})
		log.Debug().Hex("data", data).Msg("obtain")
	}()

	err := c.hold(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{true, false}),
		Payload:    []byte{0, 1, 2, 3},
	})
	assert.Nil(t, err)
	err = c.hold(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{false, true}),
		Payload:    []byte{4, 5},
	})
	assert.Nil(t, err)
	time.Sleep(time.Millisecond * 10)

}

func TestOptionalConsistentlyChannel(t *testing.T) {
	resp := make(chan *pproto.Packet)
	c := newIncomingChannel(pproto.PacketMode_PM_OPTIONAL_CONSISTENTLY, resp)
	go func() {
		for {
			r := <-resp
			log.Debug().Interface("pck", r).Msg("response")
		}
	}()

	go func() {
		data := <-c.C
		assert.EqualValues(t, data, []byte{10, 11})
		log.Debug().Hex("data", data).Msg("received")
	}()

	err := c.hold(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{true, false}),
		Payload:    []byte{0, 1, 2, 3},
	})
	assert.Nil(t, err)

	err = c.hold(&pproto.Packet{
		Id:         2,
		PartsCount: 1,
		Parts:      booleansToBytes([]bool{true}),
		Payload:    []byte{10, 11},
	})
	time.Sleep(time.Millisecond * 10)
	assert.Nil(t, err)
	err = c.hold(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{false, true}),
		Payload:    []byte{4, 5},
	})
	assert.Nil(t, err)

	err = c.hold(&pproto.Packet{
		Id:         2,
		PartsCount: 1,
		Parts:      booleansToBytes([]bool{true}),
		Payload:    []byte{10, 11},
	})
	assert.Nil(t, err)

	time.Sleep(time.Millisecond * 10)
}

func TestMandatoryChannel(t *testing.T) {
	resp := make(chan *pproto.Packet)
	c := newIncomingChannel(pproto.PacketMode_PM_MANDATORY_CONSISTENTLY, resp)
	go func() {
		for {
			r := <-resp
			log.Debug().Interface("pck", r).Msg("response")
		}
	}()

	go func() {
		data := <-c.C
		log.Debug().Hex("data", data).Msg("received")
		assert.Equal(t, data, []byte{10, 11})
		data = <-c.C
		log.Debug().Hex("data", data).Msg("received")
		assert.Equal(t, data, []byte{0, 1, 2, 3, 4, 5})

	}()

	err := c.hold(&pproto.Packet{
		Id:         2,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{false, true}),
		Payload:    []byte{4, 5},
	})
	assert.Nil(t, err)
	err = c.hold(&pproto.Packet{
		Id:         2,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{true, false}),
		Payload:    []byte{0, 1, 2, 3},
	})
	assert.Nil(t, err)
	time.Sleep(time.Millisecond * 10)
	err = c.hold(&pproto.Packet{
		Id:         1,
		PartsCount: 1,
		Parts:      booleansToBytes([]bool{true}),
		Payload:    []byte{10, 11},
	})
	assert.Nil(t, err)
	time.Sleep(time.Millisecond * 20)
}

//
//func TestMandatoryChannel(t *testing.T) {
//	retryChan := make(chan *pproto.Packet)
//	c := newMessageMandatoryChannel(retryChan)
//	packetRetryTTL = time.Millisecond * 30
//
//	go func() {
//		data := <-c.C
//		assert.EqualValues(t, data, []byte{0, 1, 2, 3, 4, 5})
//		log.Debug().Hex("data", data).Msg("received")
//	}()
//
//	go func() {
//		for i := 0; i < 4; i++ {
//			retryOrDonePacket := <-retryChan
//			log.Debug().Str("tag", retryOrDonePacket.Tag.String()).Bools("pck", bytesToBooleans(retryOrDonePacket.Parts)[:retryOrDonePacket.PartsCount]).Msg("retry")
//		}
//	}()
//	time.Sleep(time.Millisecond * 1)
//	err := c.add(&pproto.Packet{
//		Id:         1,
//		PartsCount: 3,
//		Parts:      booleansToBytes([]bool{true, false, false}),
//		Data:       []byte{0, 1},
//	})
//	assert.Nil(t, err)
//	time.Sleep(time.Millisecond * 10)
//	err = c.add(&pproto.Packet{
//		Id:         1,
//		PartsCount: 3,
//		Parts:      booleansToBytes([]bool{false, true, false}),
//		Data:       []byte{2, 3},
//	})
//	assert.Nil(t, err)
//	time.Sleep(time.Millisecond * 40)
//	err = c.add(&pproto.Packet{
//		Id:         1,
//		PartsCount: 3,
//		Parts:      booleansToBytes([]bool{false, false, true}),
//		Data:       []byte{4, 5},
//	})
//	time.Sleep(time.Millisecond * 10)
//}
