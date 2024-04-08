package tcpio

import (
	"sync"
)

type Pool struct {
	queueLen int
	retry    int
	router   map[string]*Queue
	lock     sync.Mutex
	isClose  bool
}

type PoolOptions struct {
	QueueLen int
	Retry    int
}

// NewPool create a new connection pool
func NewPool(options ...PoolOptions) *Pool {
	queueLen := 100
	retry := 3
	for _, v := range options {
		if v.QueueLen > 0 {
			queueLen = v.QueueLen
		}
		if v.Retry > 0 {
			retry = v.Retry
		}
	}
	return &Pool{
		router:   make(map[string]*Queue),
		queueLen: queueLen,
		retry:    retry,
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
