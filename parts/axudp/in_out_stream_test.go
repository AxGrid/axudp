package axudp

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"github.com/rs/zerolog/log"
	"testing"
	"time"
)

func TestInOutConnect(t *testing.T) {
	S := make(chan []byte)         // Service Channel
	R := make(chan *pproto.Packet) // ResponseChannel
	O := make(chan *pproto.Packet) // Data channel

	go func() {
		for {
			p := <-S
			log.Info().Hex("hex", p).Msg("service receive")
		}
	}()

	in := NewInStream(R, S)
	out := NewOutStream(O)

	go func() {
		for {
			select {
			case send := <-O:
				//log.Info().Interface("pck", send).Msg("send := <-O")
				in.receive(send)
			case resp := <-R:
				//log.Info().Interface("pck", resp).Msg("resp := <-R")
				out.response(resp)
			}
		}
	}()

	out.send([]byte{3, 4, 5}, pproto.PacketMode_PM_OPTIONAL)
	time.Sleep(time.Millisecond * 30)
}
