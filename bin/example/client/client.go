package main

import (
	pproto "github.com/axgrid/axudp/generated-sources/proto/axudp"
	"github.com/axgrid/axudp/transport"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"sync"
	"time"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05,000"}).Level(zerolog.DebugLevel)
	Send()
	//wa := sync.WaitGroup{}
	//for i := 0; i < 200; i++ {
	//	wa.Add(1)
	//	go func() {
	//		Send()
	//		wa.Done()
	//	}()
	//}
	//wa.Wait()
}

func Send() {
	cli, err := transport.NewClient("127.0.0.1", 9100)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to create server")
	}
	wa := sync.WaitGroup{}
	cli.ServiceListener = func(payload []byte, addr string, con transport.IConnect) error {
		t := time.Now()
		err := t.UnmarshalBinary(payload)
		if err != nil {
			return err
		}
		log.Info().Int64("ping", time.Now().Sub(t).Microseconds()).Msg("late")
		wa.Done()
		return nil
	}

	cli.Start()
	wa.Add(1)
	for i := 0; i < 500; i++ {
		bytes, err := time.Now().MarshalBinary()
		if err != nil {
			log.Fatal().Err(err).Msg("fail marshal time")
		}
		wa.Add(1)
		if i == 5 {
			cli.Send(pproto.PacketMode_PM_MANDATORY, []byte{1})
		} else {
			cli.Send(pproto.PacketMode_PM_MANDATORY, bytes)
		}
		time.Sleep(time.Millisecond * 30)
	}
	wa.Wait()
	cli.Stop()
}
