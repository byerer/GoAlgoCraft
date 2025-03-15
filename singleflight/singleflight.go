package singleflight

import "sync"

type sf struct {
	mu  sync.Mutex
	res map[string]*call
}

type call struct {
	dups int
	val  any
	wg   sync.WaitGroup
	err  error
}

func (s *sf) Do(key string, fn func() (res any, err error)) (res any, err error, shared bool) {
	s.mu.Lock()
	if s.res == nil {
		s.res = make(map[string]*call)
	}
	if c, e := s.res[key]; e {
		c.dups++
		s.mu.Unlock()

		c.wg.Wait()
		return c.val, c.err, true
	}
	c := new(call)
	c.wg.Add(1)
	s.res[key] = c
	s.mu.Unlock()
	c.val, c.err = fn()
	c.wg.Done()
	return c.val, c.err, c.dups > 0
}
