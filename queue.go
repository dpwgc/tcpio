package tcpio

import (
	"errors"
)

type Queue struct {
	pool    *Pool
	address string
	channel chan *Session
}

func newQueue(pool *Pool, address string, queueLen int) *Queue {
	channel := make(chan *Session, queueLen)
	queue := &Queue{
		pool:    pool,
		address: address,
		channel: channel,
	}
	for i := 1; i <= queueLen; i++ {
		queue.channel <- newSession(queue, i)
	}
	return queue
}

func (q *Queue) close() {
	close(q.channel)
}

func (q *Queue) popSession() (*Session, error) {
	if q.pool.isClose {
		return nil, errors.New("pool has been shut down")
	}
	s := <-q.channel
	if s.isAlive {
		return s, nil
	}
	err := s.init()
	if err != nil {
		q.channel <- s
		return nil, err
	}
	return s, nil
}

func (q *Queue) putSession(s *Session) error {
	if q.pool.isClose {
		return errors.New("pool has been shut down")
	}
	q.channel <- s
	return nil
}
