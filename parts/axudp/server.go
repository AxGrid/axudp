package axudp

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"net"
)

type UDPListenerFunc func(*UDPServerPacket) error

type UDPServerPacket struct {
	Addr   *net.UDPAddr
	Packet *pproto.Packet
}

type UDPServer struct {
	addr     *net.UDPAddr
	srv      *net.UDPConn
	listener UDPListenerFunc
}

func NewUDPServer(port int, listener UDPListenerFunc) *UDPServer {
	return &UDPServer{
		addr: &net.UDPAddr{
			Port: port,
			IP:   net.ParseIP("0.0.0.0"),
		},
		listener: listener,
	}
}

func (c *UDPServer) Start() error {
	if c.srv != nil {
		return errors.New("server already started")
	}
	ser, err := net.ListenUDP("udp", c.addr)
	if err != nil {
		return err
	}
	c.srv = ser
	go c.connections()
	return nil
}

func (c *UDPServer) connections() {
	p := make([]byte, 4096)
	log.Info().Int("port", c.addr.Port).Str("ip", c.addr.IP.String()).Msg("start udp server")
	for {
		size, addr, err := c.srv.ReadFromUDP(p)
		if err != nil {
			log.Trace().Err(err).Msg("udp error")
			break
		}
		go c.processBytes(addr, p[:size])
	}
	c.Stop()
	c.srv = nil
	log.Info().Int("port", c.addr.Port).Str("ip", c.addr.IP.String()).Msg("stop udp server")
}

func (c *UDPServer) processBytes(addr *net.UDPAddr, payload []byte) {
	var err error
	packets := dataTail(payload)
	res := make([]*UDPServerPacket, len(packets))
	for i, packetBytes := range packets {
		var pck *pproto.Packet
		err = proto.Unmarshal(packetBytes, pck)
		if err != nil {
			log.Error().Err(err).Str("addr", addr.String()).Msg("fail decode protobuf")
			break
		}
		res[i] = &UDPServerPacket{
			Addr:   addr,
			Packet: pck,
		}
	}
}

func (c *UDPServer) Stop() error {
	return c.srv.Close()
}
