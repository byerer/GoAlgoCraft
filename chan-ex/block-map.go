package chan_ex

import "sync"

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

func (m *Bm) Get(key string) string {
	m.mu.Lock()
	if v, ok := m.m[key]; ok {
		m.mu.Unlock()
		<-v.ch
		return v.val
	}
	v := &value{
		val: "",
		ch:  make(chan struct{}),
	}
	m.m[key] = v
	m.mu.Unlock()
	<-v.ch
	return v.val
}

func (m *Bm) Put(key, val string) {
	m.mu.Lock()
	if v, ok := m.m[key]; ok {
		v.val = val
		select {
		case _, open := <-v.ch:
			if !open {
				// 如果通道已经关闭，就不再关闭它
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
