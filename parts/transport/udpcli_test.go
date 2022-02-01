package transport

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestMandatoryClientServer(t *testing.T) {
	startMTU = 3
	srv, err := NewServer("0.0.0.0", 9100)
	assert.Nil(t, err)
	cli, err := NewClient("127.0.0.1", 9100)
	assert.Nil(t, err)

	srv.Start()
	cli.Start()

	var clientSendQueue [][]byte
	for i := 0; i < 50; i++ {
		clientSendQueue = append(clientSendQueue, make([]byte, i+1))
	}
	var clientReceiveQueue []int
	var serverReceiveQueue []int

	srv.serviceListener = func(payload []byte, addr string, con IConnect) error {
		log.Debug().Hex("payload", payload).Str("addr", addr).Msg("server receive")
		serverReceiveQueue = append(serverReceiveQueue, len(payload))
		go func() {
			con.Send(pproto.PacketMode_PM_MANDATORY, payload)
		}()
		return nil
	}

	cli.serviceListener = func(payload []byte, addr string, con IConnect) error {
		log.Debug().Hex("payload", payload).Str("addr", addr).Msg("client receive")
		clientReceiveQueue = append(clientReceiveQueue, len(payload))
		return nil
	}
	for _, pl := range clientSendQueue {
		cli.Send(pproto.PacketMode_PM_MANDATORY, pl)
	}
	time.Sleep(time.Millisecond * 200)
	assert.Equal(t, len(clientReceiveQueue), len(clientSendQueue))

	for i := 0; i < 50; i++ {
		found := false
		for _, d := range clientReceiveQueue {
			if d == i+1 {
				found = true
				break
			}
		}
		assert.True(t, found)
	}

	log.Info().Ints("client", clientReceiveQueue).Msg("done")
	log.Info().Ints("server", serverReceiveQueue).Msg("done")

}

func TestMandatoryConsistentClientServer(t *testing.T) {
	startMTU = 3
	srv, err := NewServer("0.0.0.0", 9100)
	assert.Nil(t, err)
	cli, err := NewClient("127.0.0.1", 9100)
	assert.Nil(t, err)

	srv.Start()
	cli.Start()

	var clientSendQueue [][]byte
	for i := 0; i < 50; i++ {
		clientSendQueue = append(clientSendQueue, make([]byte, i+1))
	}
	var clientReceiveQueue []int
	var serverReceiveQueue []int

	srv.serviceListener = func(payload []byte, addr string, con IConnect) error {
		log.Debug().Hex("payload", payload).Str("addr", addr).Int("len", len(payload)).Msg("server receive")
		serverReceiveQueue = append(serverReceiveQueue, len(payload))
		//go func() {
		//	con.Send(pproto.PacketMode_PM_MANDATORY_CONSISTENTLY, payload)
		//}()
		return nil
	}

	cli.serviceListener = func(payload []byte, addr string, con IConnect) error {
		log.Debug().Hex("payload", payload).Str("addr", addr).Int("len", len(payload)).Msg("client receive")
		clientReceiveQueue = append(clientReceiveQueue, len(payload))
		return nil
	}
	for _, pl := range clientSendQueue {
		cli.Send(pproto.PacketMode_PM_MANDATORY_CONSISTENTLY, pl)
	}
	time.Sleep(time.Millisecond * 2500)
	log.Info().Ints("client", clientReceiveQueue).Msg("done")
	log.Info().Ints("server", serverReceiveQueue).Msg("done")

	assert.Equal(t, len(clientReceiveQueue), len(clientSendQueue))

	//for i := 0; i < 50; i++ {
	//	found := false
	//	for _, d := range clientReceiveQueue {
	//		if d == i+1 {
	//			found = true
	//			break
	//		}
	//	}
	//	assert.True(t, found)
	//}

}

func TestMain(m *testing.M) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05,000"}).Level(zerolog.TraceLevel)
	code := m.Run()
	os.Exit(code)
}
