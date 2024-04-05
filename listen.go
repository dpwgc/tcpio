package tcpio

import (
	"net"
)

type HandleFunc func(conn *net.TCPConn, err error)

// Listen start tcp listening
func Listen(address string, handleFunc HandleFunc) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}
	listen, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	for {
		go handleFunc(listen.AcceptTCP())
	}
}
