package tcpio

import "net"

type HandleFunc func(conn net.Conn, err error)

// Listen start tcp listening
func Listen(address string, handleFunc HandleFunc) error {
	listen, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	for {
		go handleFunc(listen.Accept())
	}
}
