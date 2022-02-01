package transport

import (
	pproto "axudp/target/generated-sources/proto/axudp"
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
)

type Client struct {
	connection      *Connection
	errorChan       chan ConnectionResponse
	close           chan bool
	serviceError    func(err error, addr string)
	serviceListener func(payload []byte, addr string, con IConnect) error
	cli             net.Conn
}

func NewClient(host string, port int) (*Client, error) {
	rs := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("udp", rs)
	if err != nil {
		return nil, err
	}
	errorChan := make(chan ConnectionResponse)
	res := &Client{
		errorChan: errorChan,
		cli:       conn,
	}

	res.connection = NewConnection(rs, errorChan, func(bytes []byte) {
		log.Trace().Str("class", "udpcli").Hex("hex", bytes).Msg("send")
		_, werr := conn.Write(addSize(bytes))
		if werr != nil {
			log.Error().Err(err).Msg("client connection error")
			errorChan <- ConnectionResponse{err: werr, address: rs}
		}
	}, func(bytes []byte) {
		if res.connection != nil {
			serr := res.serviceListener(bytes, rs, res)
			if serr != nil {
				res.errorChan <- ConnectionResponse{err: serr, address: rs}
			}
		} else {
			log.Warn().Hex("receive", bytes).Msg("client listener not set")
		}
	})
	return res, nil
}

func (c *Client) Start() {
	go c.readLoop()
	go func() {
		for {
			select {
			case err := <-c.errorChan:
				if c.serviceError != nil {
					c.serviceError(err.err, c.connection.remoteAddrStr)
				}
				c.cli.Close()
				return
			case <-c.close:
				c.cli.Close()
				return
			}
		}
	}()
}

func (c *Client) readLoop() {
	p := make([]byte, 2048)
	for {
		cnt, err := c.cli.Read(p)
		if err != nil {
			c.errorChan <- ConnectionResponse{err: err, address: c.connection.remoteAddrStr}
		}
		slice := make([]byte, cnt)
		copy(slice, p[:cnt])
		go func(data []byte) {
			log.Trace().Str("class", "udpcli").Hex("hex", slice).Msg("recv")
			for _, pck := range dataTail(data) {
				c.connection.receive(pck)
			}
		}(slice)
	}
}

func (c *Client) Send(mode pproto.PacketMode, payload []byte) {
	c.connection.Send(mode, payload)
}
