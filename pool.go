package tcpio

import "time"

type Pool struct {
	queueLen int
	retry    int
	timeout  time.Duration
	deadline time.Duration
	router   map[string]*Queue
	isClose  bool
}

type PoolOptions struct {
	QueueLen int
	Retry    int
	Timeout  time.Duration
	Deadline time.Duration
}

// NewPool create a new connection pool
func NewPool(options ...PoolOptions) *Pool {
	queueLen := 100
	retry := 3
	timeout := time.Second * 30
	deadline := time.Minute * 5
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
		if v.Deadline.Milliseconds() > 0 {
			deadline = v.Deadline
		}
	}
	return &Pool{
		router:   make(map[string]*Queue),
		queueLen: queueLen,
		retry:    retry,
		timeout:  timeout,
		deadline: deadline,
		isClose:  false,
	}
}

// Session get a alive session
func (p *Pool) Session(addr string) (*Session, error) {
	if p.router[addr] == nil {
		p.router[addr] = newQueue(p, addr, p.queueLen)
	}
	var err error = nil
	var session *Session = nil
	for i := 0; i <= p.retry; i++ {
		session, err = p.router[addr].popSession()
		if err == nil {
			break
		}
	}
	session.isFinish = false
	return session, err
}

// Close shut down connection pool
func (p *Pool) Close() {
	p.isClose = true
	for _, v := range p.router {
		v.close()
	}
}
