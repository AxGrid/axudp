package main

import (
	pproto "github.com/axgrid/axudp/generated-sources/proto/axudp"
	"github.com/axgrid/axudp/transport"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05,000"}).Level(zerolog.InfoLevel)
	cli, err := transport.NewClient("127.0.0.1", 9100)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to create server")
	}

	cli.ServiceListener = func(payload []byte, addr string, con transport.IConnect) error {
		//log.Info().Str("addr", addr).Hex("payload", payload).Msg("response")
		t := time.Now()
		err := t.UnmarshalBinary(payload)
		if err != nil {
			return err
		}
		log.Info().Int64("ping", time.Now().Sub(t).Microseconds()).Msg("late")
		return nil
	}

	cli.Start()
	for i := 0; i < 10; i++ {
		bytes, err := time.Now().MarshalBinary()
		if err != nil {
			log.Fatal().Err(err).Msg("fail marshal time")
		}
		cli.Send(pproto.PacketMode_PM_MANDATORY, bytes)
	}

	time.Sleep(time.Millisecond * 1000)

}
