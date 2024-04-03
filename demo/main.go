package main

import (
	"bufio"
	"fmt"
	"net"
	"tcpio"
	"time"
)

func main() {

	go server()

	time.Sleep(1 * time.Second)

	pool := tcpio.NewPool(tcpio.PoolOptions{
		QueueLen: 1,
	})

	err := useSession(pool, "0.0.0.0:8081", "hello world session 1")
	if err != nil {
		panic(err)
	}

	err = useSession(pool, "0.0.0.0:8081", "hello world session 2")
	if err != nil {
		panic(err)
	}
}

func useSession(pool *tcpio.Pool, addr, text string) error {
	session, err := pool.Session(addr)
	if err != nil {
		return err
	}
	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("session id:", session.ID(), "| write:", text)
		_, err = session.Write([]byte(text))
		if err != nil {
			return err
		}
	}
	err = session.Finish()
	if err != nil {
		return err
	}
	return nil
}

func server() {
	err := tcpio.Listen("0.0.0.0:8081", handleFunc)
	if err != nil {
		panic(err)
	}
}

func handleFunc(conn net.Conn, err error) {
	if err != nil {
		panic(err)
	}
	for {
		reader := bufio.NewReader(conn)
		var buf [1024]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			panic(err)
		}
		fmt.Println("server read:", string(buf[:n]))
	}
}
