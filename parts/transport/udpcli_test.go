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

func TestClientServer(t *testing.T) {
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

	srv.serviceListener = func(payload []byte, addr string, con IConnect) error {
		log.Debug().Hex("payload", payload).Str("addr", addr).Msg("server receive")
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
				log.Debug().Int("d", d).Msg("found")
				found = true
				break
			}
		}
		assert.True(t, found)
	}

}

func TestMain(m *testing.M) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "15:04:05,000"}).Level(zerolog.DebugLevel)
	code := m.Run()
	os.Exit(code)
}
