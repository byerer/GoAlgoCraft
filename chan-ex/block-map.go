package chan_ex

import (
	"errors"
	"sync"
	"time"
)

type Bm struct {
	mu sync.Mutex
	m  map[string]*value
}

func NewBm() *Bm {
	return &Bm{
		mu: sync.Mutex{},
		m:  map[string]*value{},
	}
}

type value struct {
	val string
	ch  chan struct{}
}

func (m *Bm) Get(key string, maxTime time.Duration) (string, error) {
	m.mu.Lock()
	if v, ok := m.m[key]; ok {
		m.mu.Unlock()
		select {
		case <-time.After(maxTime):
			return "", errors.New("timeout")
		case <-v.ch:
			return v.val, nil
		}
	}
	v := &value{
		val: "",
		ch:  make(chan struct{}),
	}
	m.m[key] = v
	m.mu.Unlock()
	select {
	case <-time.After(maxTime):
		return "", errors.New("timeout")
	case <-v.ch:
		return v.val, nil
	}
}

func (m *Bm) Put(key, val string) {
	m.mu.Lock()
	if v, ok := m.m[key]; ok {
		v.val = val
		select {
		case _, open := <-v.ch:
			if !open {
				m.mu.Unlock()
				return
			}
		default:
			// 通道未关闭，关闭它
			close(v.ch)
		}
		m.mu.Unlock()
		return
	}
	v := &value{
		val: val,
		ch:  make(chan struct{}),
	}
	m.m[key] = v
	m.mu.Unlock()
	close(v.ch)
	return
}
