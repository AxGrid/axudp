package transport

import (
	pproto "github.com/axgrid/axudp/generated-sources/proto/axudp"
	udpchan2 "github.com/axgrid/axudp/udpchan"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"net"
)

type Connection struct {
	remoteAddr             net.UDPAddr
	remoteAddrStr          string
	close                  chan bool
	mtu                    int
	errorChan              chan error
	inMandatory            *udpchan2.InMandatoryChannel
	inMandatoryConsistent  *udpchan2.InMandatoryChannel
	outMandatory           *udpchan2.OutMandatoryChannel
	outMandatoryConsistent *udpchan2.OutMandatoryConsistentChannel
	log                    zerolog.Logger
}

func NewConnection(remoteAddrStr string, errorServerChan chan ConnectionResponse, send func([]byte), service func([]byte)) *Connection {
	res := &Connection{
		remoteAddrStr: remoteAddrStr,
		mtu:           startMTU,
		errorChan:     make(chan error),
		close:         make(chan bool, 1),
		log:           log.With().Str("addr", remoteAddrStr).Logger(),
	}
	servChan := make(chan []byte)
	sendChan := make(chan []byte)
	res.inMandatory = udpchan2.NewInMandatoryChannel(pproto.PacketMode_PM_MANDATORY, res.mtu, servChan, sendChan, res.errorChan, res.log)
	res.inMandatoryConsistent = udpchan2.NewInMandatoryChannel(pproto.PacketMode_PM_MANDATORY_CONSISTENTLY, res.mtu, servChan, sendChan, res.errorChan, res.log)
	res.outMandatory = udpchan2.NewOutMandatoryChannel(res.mtu, sendChan, res.errorChan, res.log)
	res.outMandatoryConsistent = udpchan2.NewOutMandatoryConsistentChannel(res.mtu, sendChan, res.errorChan, res.log)

	go func() {
		for {
			select {
			case <-res.close:
				res.stopAllChannels()
				return
			case b := <-sendChan:
				send(b)
			case err := <-res.errorChan:
				errorServerChan <- ConnectionResponse{err: err, address: res.remoteAddrStr}
				res.stopAllChannels()
				return
			case payload := <-servChan:
				service(payload)
			}
		}
	}()
	return res
}

func (c *Connection) stopAllChannels() {
	c.outMandatory.Stop()
	c.inMandatory.Stop()
	c.outMandatoryConsistent.Stop()
	c.inMandatoryConsistent.Stop()
}

func (c *Connection) receive(data []byte) {
	var pck pproto.Packet
	err := proto.Unmarshal(data, &pck)
	if err != nil {
		c.errorChan <- err
		return
	}

	switch pck.Tag {
	case pproto.PacketTag_PT_DONE:
		switch pck.Mode {
		case pproto.PacketMode_PM_MANDATORY:
			c.outMandatory.Response(&pck)
		case pproto.PacketMode_PM_MANDATORY_CONSISTENTLY:
			c.outMandatoryConsistent.Response(&pck)
		}
	case pproto.PacketTag_PT_PAYLOAD, pproto.PacketTag_PT_GZIP_PAYLOAD:
		switch pck.Mode {
		case pproto.PacketMode_PM_MANDATORY:
			log.Trace().Bools("parts", bytesToBooleans(pck.Parts)[:pck.PartsCount]).Uint64("id", pck.Id).Msg("connection receive")
			c.inMandatory.Receive(&pck)
		case pproto.PacketMode_PM_MANDATORY_CONSISTENTLY:
			log.Trace().Bools("parts", bytesToBooleans(pck.Parts)[:pck.PartsCount]).Uint64("id", pck.Id).Msg("connection receive")
			c.inMandatoryConsistent.Receive(&pck)

		}
	}

}

func (c *Connection) Send(mode pproto.PacketMode, payload []byte) {
	switch mode {
	case pproto.PacketMode_PM_MANDATORY:
		c.outMandatory.Send(payload)
	case pproto.PacketMode_PM_MANDATORY_CONSISTENTLY:
		c.outMandatoryConsistent.Send(payload)
	default:
		log.Fatal().Msgf("mode %s not implemented", mode.String())
	}
}
