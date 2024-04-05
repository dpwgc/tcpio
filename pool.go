package tcpio

import (
	"sync"
	"time"
)

type Pool struct {
	queueLen int
	retry    int
	timeout  time.Duration
	router   map[string]*Queue
	lock     sync.Mutex
	isClose  bool
}

type PoolOptions struct {
	QueueLen int
	Retry    int
	Timeout  time.Duration
}

// NewPool create a new connection pool
func NewPool(options ...PoolOptions) *Pool {
	queueLen := 100
	retry := 3
	timeout := time.Minute * 5
	for _, v := range options {
		if v.QueueLen > 0 {
			queueLen = v.QueueLen
		}
		if v.Retry > 0 {
			retry = v.Retry
		}
		if v.Timeout.Milliseconds() > 0 {
			timeout = v.Timeout
		}
	}
	return &Pool{
		router:   make(map[string]*Queue),
		queueLen: queueLen,
		retry:    retry,
		timeout:  timeout,
		lock:     sync.Mutex{},
		isClose:  false,
	}
}

// Session get a alive session
func (p *Pool) Session(addr string) (*Session, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if p.router[addr] == nil {
		p.router[addr] = newQueue(p, addr, p.queueLen)
	}
	session, err := p.router[addr].popSession()
	if err != nil {
		return nil, err
	}
	session.isFree = false
	return session, nil
}

// Close shut down connection pool
func (p *Pool) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.isClose = true
	for _, v := range p.router {
		v.close()
	}
}
