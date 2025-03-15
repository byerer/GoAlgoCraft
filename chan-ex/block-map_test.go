package chan_ex

import (
	"fmt"
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
		val := m.Get("111")
		fmt.Println(val)
		if val != "222" {
			t.Errorf("expect 222, but got %s", val)
			return
		}
	}()

	go func() {
		defer wg.Done()
		val := m.Get("111")
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
		val := m.Get("111")
		if val != "222" {
			t.Errorf("expect 222, but got %s", val)
			return
		}
	}()
	wg.Wait()
}
