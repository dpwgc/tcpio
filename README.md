# TCPIO

## A simple TCP connection pool

***

## How to use

### client example

#### tcp connection pool

```go
// init a new connection pool (global object)
var pool = tcpio.NewPool()

// client example
func clientExample() {

	// get a session
	session, _ := pool.Session("0.0.0.0:8081")

	// write request message
	_, _ = session.Write([]byte("hello world"))

	// read response message
	var response [1024]byte
	n, _ := session.Read(response[:])

	// print response message
	fmt.Println("response:", string(response[:n]))

	// free the session
	_ = session.Finish()
}
```

### server example

#### tcp listening

```go
// server example
func serverExample() {

	// start tcp listening
	_ = tcpio.Listen("0.0.0.0:8081", func(conn net.Conn, err error) {

		// read request message
		var request [1024]byte
		n, _ := conn.Read(request[:])

		// print request message
		fmt.Println("request:", string(request[:n]))

		// write response message
		_, _ = conn.Write([]byte("hi"))
	})
}
```

***

## Function
* tcpio
  * `NewPool` create a new connection pool
  * `Listen` start tcp listening
* pool
  * `Session` get a alive session
  * `Close` shut down connection pool
* session
  * `Write` tcp write
  * `Read` tcp read
  * `Finish` free the session