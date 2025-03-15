package chan_ex

import (
	"sync"
	"testing"
	"time"
)

func TestBM(t *testing.T) {
	m := NewBm()
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		defer wg.Done()
		val, _ := m.Get("111", 100*time.Millisecond)
		if val != "222" {
			t.Errorf("expect 222, but got %s", val)
			return
		}
	}()

	go func() {
		defer wg.Done()
		val, _ := m.Get("111", 100*time.Millisecond)
		if val != "222" {
			t.Errorf("expect 222, but got %s", val)
			return
		}
	}()

	go func() {
		defer wg.Done()
		time.Sleep(10 * time.Millisecond)
		m.Put("111", "222")
	}()

	go func() {
		defer wg.Done()
		val, _ := m.Get("111", 100*time.Millisecond)
		if val != "222" {
			t.Errorf("expect 222, but got %s", val)
			return
		}
	}()
	wg.Wait()
}

func TestTimeout(t *testing.T) {
	m := NewBm()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := m.Get("111", 10*time.Millisecond)
		if err == nil {
			t.Errorf("expect timeout, but got nil")
			return
		}
	}()

	go func() {
		defer wg.Done()
		time.Sleep(100 * time.Millisecond)
		m.Put("111", "222")
	}()
	wg.Wait()
}
