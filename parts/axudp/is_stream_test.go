package axudp

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"github.com/rs/zerolog/log"
	"testing"
	"time"
)

func TestInOptional(t *testing.T) {
	R := make(chan *pproto.Packet)
	C := make(chan []byte)
	go func() {
		for {
			p := <-R
			log.Info().Uint64("id", p.Id).Bools("parts", bytesToBooleans(p.Parts)[:p.PartsCount]).Msg("response")
		}
	}()

	c := newInOptionalStream(R, C)

	go func() {
		for {
			p := <-c.C
			log.Info().Hex("payload", p).Msg("recv")
		}
	}()

	waitDuplicatePacket = time.Millisecond * 10

	c.receive(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{false, true}),
		Payload:    []byte{4, 5},
	})

	c.receive(&pproto.Packet{
		Id:         2,
		PartsCount: 1,
		Parts:      booleansToBytes([]bool{true}),
		Payload:    []byte{8, 9},
	})

	time.Sleep(time.Millisecond * 1)

	c.receive(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{false, true}),
		Payload:    []byte{4, 5},
	})

	c.receive(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{true, false}),
		Payload:    []byte{1, 2, 3},
	})

	c.receive(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{true, false}),
		Payload:    []byte{1, 2, 3},
	})

	time.Sleep(time.Millisecond * 20)

	c.receive(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{false, true}),
		Payload:    []byte{4, 5},
	})

	c.receive(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{true, false}),
		Payload:    []byte{1, 2, 3},
	})

	time.Sleep(time.Millisecond * 10)

}

func TestInOptionalConsistently(t *testing.T) {
	R := make(chan *pproto.Packet)
	C := make(chan []byte)
	go func() {
		for {
			p := <-R
			log.Info().Uint64("id", p.Id).Bools("parts", bytesToBooleans(p.Parts)[:p.PartsCount]).Msg("response")
		}
	}()

	c := newInOptionalConsistentlyStream(R, C)

	go func() {
		for {
			p := <-c.C
			log.Info().Hex("payload", p).Msg("recv")
		}
	}()

	waitDuplicatePacket = time.Millisecond * 10

	c.receive(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{false, true}),
		Payload:    []byte{4, 5},
	})

	c.receive(&pproto.Packet{
		Id:         2,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{false, true}),
		Payload:    []byte{8, 9},
	})

	time.Sleep(time.Millisecond * 1)

	c.receive(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{false, true}),
		Payload:    []byte{4, 5},
	})

	c.receive(&pproto.Packet{
		Id:         2,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{true, false}),
		Payload:    []byte{6, 7},
	})

	c.receive(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{true, false}),
		Payload:    []byte{1, 2, 3},
	})

	c.receive(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{true, false}),
		Payload:    []byte{1, 2, 3},
	})

	time.Sleep(time.Millisecond * 20)

	c.receive(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{false, true}),
		Payload:    []byte{4, 5},
	})

	c.receive(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{true, false}),
		Payload:    []byte{1, 2, 3},
	})

	time.Sleep(time.Millisecond * 10)
}

func TestInMandatory(t *testing.T) {
	R := make(chan *pproto.Packet)
	C := make(chan []byte)
	go func() {
		for {
			p := <-R
			log.Info().Uint64("id", p.Id).Bools("parts", bytesToBooleans(p.Parts)[:p.PartsCount]).Msg("response")
		}
	}()

	c := newInMandatoryStream(R, pproto.PacketMode_PM_MANDATORY, C)

	go func() {
		for {
			p := <-c.C
			log.Info().Hex("payload", p).Msg("recv")
		}
	}()

	waitDuplicatePacket = time.Millisecond * 10

	c.receive(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{false, true}),
		Payload:    []byte{4, 5},
	})

	c.receive(&pproto.Packet{
		Id:         2,
		PartsCount: 1,
		Parts:      booleansToBytes([]bool{true}),
		Payload:    []byte{8, 9},
	})

	time.Sleep(time.Millisecond * 1)

	c.receive(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{false, true}),
		Payload:    []byte{4, 5},
	})

	c.receive(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{true, false}),
		Payload:    []byte{1, 2, 3},
	})

	c.receive(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{true, false}),
		Payload:    []byte{1, 2, 3},
	})

	time.Sleep(time.Millisecond * 20)

	c.receive(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{false, true}),
		Payload:    []byte{4, 5},
	})

	c.receive(&pproto.Packet{
		Id:         1,
		PartsCount: 2,
		Parts:      booleansToBytes([]bool{true, false}),
		Payload:    []byte{1, 2, 3},
	})

	time.Sleep(time.Millisecond * 10)
}
