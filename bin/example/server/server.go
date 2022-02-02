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
	srv, err := transport.NewServer("0.0.0.0", 9100)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to create server")
	}
	srv.ServiceListener = func(payload []byte, addr string, con transport.IConnect) error {
		log.Info().Str("addr", addr).Hex("payload", payload).Msg("receive")
		go func() {
			con.Send(pproto.PacketMode_PM_MANDATORY, payload)
		}()
		return nil
	}

	srv.Start()
	time.Sleep(time.Hour * 100)
}
