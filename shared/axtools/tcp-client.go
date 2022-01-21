package axtools

import (
	"bufio"
	"encoding/binary"
	"github.com/rs/zerolog/log"
	"net"
)

func TcpClient(address string, bufferSize int, handle func(out []byte) error) (func(in []byte) error, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	log.Info().Interface("conn", conn).Msg("create client")
	if err != nil {
		return nil, err
	}
	go func() {
	MainLoop:
		for {
			buffer := make([]byte, bufferSize)
			readSize, err := bufio.NewReader(conn).Read(buffer)
			if err != nil {
				break MainLoop
			}
			buffer = buffer[:readSize]
			size := int(binary.BigEndian.Uint16(buffer[:2]))
			for len(buffer) >= size {
				pck := buffer[2 : size+2]
				log.Info().Bytes("pck", pck).Msg("<-")
				buffer = buffer[size+2:]
				if err != nil {
					log.Error().Err(err).Msg("protobuf error")
					break MainLoop
				}
				if handle(pck) != nil {
					break MainLoop
				}
				size = int(binary.BigEndian.Uint16(buffer[:2]))
			}
		}
		log.Error().Msg("connection lost")
	}()
	return func(in []byte) (err error) {
		if err != nil {
			return
		}
		sizeBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(sizeBytes, uint16(len(in)))
		log.Info().Bytes("pck", in).Msg("->")
		_, err = conn.Write(append(sizeBytes, in...))
		if err != nil {
			return
		}
		return
	}, nil
}
