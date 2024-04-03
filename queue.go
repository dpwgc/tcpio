package tcpio

import (
	"errors"
	"net"
	"time"
)

type Queue struct {
	pool    *Pool
	addr    string
	channel chan *Session
	isClose bool
}

func newQueue(pool *Pool, addr string, queueLen int) *Queue {
	channel := make(chan *Session, queueLen)
	queue := &Queue{
		pool:    pool,
		addr:    addr,
		channel: channel,
		isClose: false,
	}
	for i := 0; i < queueLen; i++ {
		queue.channel <- queue.newSession()
	}
	return queue
}

func (q *Queue) close() {
	q.isClose = true
	close(q.channel)
}

func (q *Queue) popSession() (*Session, error) {
	if q.isClose {
		return nil, errors.New("session pool has been shut down")
	}
	s := <-q.channel
	if s.isAlive {
		return s, nil
	}
	tcpConn, err := net.DialTimeout("tcp", q.addr, q.pool.timeout)
	if err != nil {
		q.channel <- s
		return nil, err
	}
	err = tcpConn.SetDeadline(time.Now().Add(q.pool.deadline))
	if err != nil {
		q.channel <- s
		return nil, err
	}
	s.isAlive = true
	s.conn = tcpConn
	return s, nil
}

func (q *Queue) putSession(s *Session) error {
	if q.isClose {
		return errors.New("session pool has been shut down")
	}
	q.channel <- s
	return nil
}

func (q *Queue) newSession() *Session {
	return newSession(q)
}
