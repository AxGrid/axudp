package transport

import (
	pproto "github.com/axgrid/axudp/generated-sources/proto/axudp"
	"github.com/rs/zerolog/log"
	"net"
	"sync"
)

type ConnectionResponse struct {
	err     error
	address string
}

var startMTU = 400

type IConnect interface {
	Send(pproto.PacketMode, []byte)
}

type Server struct {
	connections     map[string]*Connection
	errorChan       chan ConnectionResponse
	serviceError    func(err error, addr string)
	serviceListener func(payload []byte, addr string, con IConnect) error
	srv             *net.UDPConn
	lock            sync.RWMutex
	close           chan bool
}

func (s *Server) Send(data []byte, addr *net.UDPAddr) error {
	_, err := s.srv.WriteToUDP(addSize(data), addr)
	return err
}

func NewServer(host string, port int) (*Server, error) {
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(host),
	}
	srv, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return nil, err
	}
	res := &Server{
		srv:         srv,
		close:       make(chan bool, 1),
		errorChan:   make(chan ConnectionResponse),
		connections: map[string]*Connection{},
	}
	return res, nil
}

func (s *Server) Start() {
	go s.readLoop()
	go func() {
		for {
			select {
			case <-s.close:
				s.srv.Close()
				return
			case err := <-s.errorChan:
				log.Error().Err(err.err).Str("addr", err.address).Msg("error in connection")
				s.lock.Lock()
				conn, ok := s.connections[err.address]
				if ok && s.serviceError != nil {
					s.serviceError(err.err, conn.remoteAddrStr)
				}
				delete(s.connections, err.address)
				s.lock.Unlock()
			}
		}
	}()
}

func (s *Server) Stop() {
	s.close <- true
}

func (s *Server) readLoop() {
	p := make([]byte, 2048)
	for {
		cnt, remoteAddr, err := s.srv.ReadFromUDP(p)
		if err != nil {
			log.Error().Err(err).Str("remote", remoteAddr.String()).Msg("error read from udp")
			break
		}
		slice := make([]byte, cnt)
		copy(slice, p[:cnt])
		go func(data []byte, remoteAddr *net.UDPAddr) {
			log.Trace().Str("class", "udpsrv").Hex("hex", data).Msg("recv thread")
			s.lock.Lock()
			conn, ok := s.connections[remoteAddr.String()]
			if !ok {
				conn = NewConnection(remoteAddr.String(), s.errorChan, func(bytes []byte) {
					log.Trace().Str("class", "udpsrv").Hex("hex", p[:cnt]).Msg("send")
					if cerr := s.Send(bytes, remoteAddr); cerr != nil {
						s.errorChan <- ConnectionResponse{address: conn.remoteAddrStr, err: cerr}
					}
				}, func(bytes []byte) {
					if s.serviceListener != nil {
						serr := s.serviceListener(bytes, conn.remoteAddrStr, conn)
						if serr != nil {
							s.errorChan <- ConnectionResponse{address: conn.remoteAddrStr, err: serr}
						}
					} else {
						log.Warn().Hex("receive", bytes).Msg("server listener not set")
					}
				})
				s.connections[remoteAddr.String()] = conn
			}
			s.lock.Unlock()
			for _, pck := range dataTail(data) {
				log.Trace().Hex("pck", pck).Str("addr", remoteAddr.String()).Msg("into connection")
				conn.receive(pck)
			}
		}(slice, remoteAddr)
	}
	log.Debug().Msg("server read loop closed")
}
