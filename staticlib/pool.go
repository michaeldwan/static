package staticlib

import (
	"runtime"
	"sync"
)

type Pool struct {
	wg    *sync.WaitGroup
	ready chan bool
	m     sync.Mutex
	stop  bool
}

type JobFunc func()

func NewPool(concurrencyLimit int) *Pool {
	if concurrencyLimit < 1 {
		concurrencyLimit = runtime.NumCPU()
	}

	pool := &Pool{
		wg:    &sync.WaitGroup{},
		ready: make(chan bool, concurrencyLimit),
	}

	for i := 0; i < concurrencyLimit; i++ {
		pool.ready <- true
	}

	return pool
}

func (p *Pool) Run(job JobFunc) {
	<-p.ready
	p.wg.Add(1)

	go func() {
		job()
		p.ready <- true
		p.wg.Done()
	}()
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Lock() {
	p.m.Lock()
}

func (p *Pool) Unlock() {
	p.m.Unlock()
}
