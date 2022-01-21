package axudp

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func TestCheckPacketSize(t *testing.T) {
	pck := &pproto.Packet{
		Payload:    []byte{},
		Parts:      booleansToBytes([]bool{false, false, true, true, false, true, true}),
		PartsCount: 12,
		Mode:       pproto.PacketMode_PM_OPTIONAL,
		Tag:        pproto.PacketTag_PT_GZIP_PAYLOAD,
	}
	bytes, err := proto.Marshal(pck)
	assert.Nil(t, err)
	log.Info().Int("size", len(bytes)).Msg("proto size")
}

func TestStartServer(t *testing.T) {
	srv := NewUDPServer(100, func(packet *UDPServerPacket) error {
		return nil
	})
	assert.Nil(t, srv.Start())
	time.Sleep(time.Millisecond * 50)
	assert.Nil(t, srv.Stop())
}
