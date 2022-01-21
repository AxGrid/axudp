package axudp

type UDPClient struct {
	close chan bool
}

func NewUDPClient(addr string) *UDPClient {
	res := &UDPClient{
		close: make(chan bool, 1),
	}
	//conn, err := net.Dial("udp", addr)
	//go func() {
	//	for {
	//		conn.RE
	//	}
	//}()
	return res
}
