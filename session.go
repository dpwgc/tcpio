package tcpio

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type Session struct {
	id       string
	queue    *Queue
	conn     net.Conn
	isAlive  bool
	isFinish bool
}

func newSession(queue *Queue) *Session {
	return &Session{
		id:       fmt.Sprintf("%s-%s-%v", queue.addr, uuidGen(), time.Now().UnixMilli()),
		queue:    queue,
		conn:     nil,
		isAlive:  false,
		isFinish: false,
	}
}

func (s *Session) ID() string {
	return s.id
}

// Read tcp read
func (s *Session) Read(buf []byte) (int, error) {
	if s.isFinish {
		return 0, errors.New("session is complete")
	}
	var n = 0
	var err error = nil
	for i := 0; i <= s.queue.pool.retry; i++ {
		n, err = s.read(buf)
		if err == nil {
			break
		}
	}
	return n, err
}

// Write tcp read
func (s *Session) Write(buf []byte) (int, error) {
	if s.isFinish {
		return 0, errors.New("session is complete")
	}
	var n = 0
	var err error = nil
	for i := 0; i <= s.queue.pool.retry; i++ {
		n, err = s.write(buf)
		if err == nil {
			break
		}
	}
	return n, err
}

func (s *Session) RemoteAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *Session) LocalAddr() net.Addr {
	return s.conn.LocalAddr()
}

// Finish free the session
func (s *Session) Finish() error {
	if s.isFinish {
		return errors.New("session is complete")
	}
	s.isFinish = true
	return s.queue.putSession(s)
}

func (s *Session) close() {
	_ = s.conn.Close()
	s.isAlive = false
}

func (s *Session) read(buf []byte) (int, error) {
	err := s.conn.SetDeadline(time.Now().Add(s.queue.pool.deadline))
	if err != nil {
		s.close()
		return 0, err
	}
	n, err := s.conn.Read(buf)
	if err != nil {
		s.close()
		return 0, err
	}
	return n, nil
}

func (s *Session) write(buf []byte) (int, error) {
	err := s.conn.SetDeadline(time.Now().Add(s.queue.pool.deadline))
	if err != nil {
		s.close()
		return 0, err
	}
	n, err := s.conn.Write(buf)
	if err != nil {
		s.close()
		return 0, err
	}
	return n, nil
}
