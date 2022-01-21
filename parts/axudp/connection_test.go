package axudp

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddressEquality(t *testing.T) {
	a := net.UDPAddr{
		Port: 7100,
		IP:   net.IPv4(127, 0, 0, 1),
		Zone: "test",
	}
	b := net.UDPAddr{
		Port: 7100,
		IP:   net.IPv4(127, 0, 0, 1),
		Zone: "test",
	}
	c := net.UDPAddr{
		Port: 7100,
		IP:   net.IPv4(127, 0, 0, 2),
		Zone: "test",
	}
	assert.Equal(t, a, b)
	assert.NotEqual(t, a, c)
}

func TestNewUDPServer(t *testing.T) {
	srv := NewUDPServer(4100, addOrUpdateConnection)
	assert.Nil(t, srv.Start())
	time.Sleep(time.Millisecond * 50)
	assert.Nil(t, srv.Stop())
}
