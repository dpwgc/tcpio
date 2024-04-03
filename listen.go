package tcpio

import "net"

type HandleFunc func(conn net.Conn, err error)

// Listen start tcp listening
func Listen(addr string, handleFunc HandleFunc) error {
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	for {
		go handleFunc(listen.Accept())
	}
}
