package chan_ex

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/sync/singleflight"
)

func TestSingleflight(t *testing.T) {
	sf := &singleflight.Group{}

	t.Run("Test_Normal_Shared", func(t *testing.T) {
		var res1, res2 any
		var err1, err2 error
		var shared1, shared2 bool

		// 第一次调用
		res1, err1, shared1 = sf.Do("key1", func() (any, error) {
			fmt.Println("Running first function...")
			time.Sleep(1 * time.Second)
			return "result1", nil
		})

		// 第二次调用相同的key，应该共享结果
		res2, err2, shared2 = sf.Do("key1", func() (any, error) {
			return nil, nil // 这里不应该执行到
		})

		// 第一次请求的检查
		if err1 != nil || shared1 != false || res1 != "result1" {
			t.Errorf("unexpected result from first call: %v, %v, %v", res1, err1, shared1)
		}

		// 第二次请求的检查
		if err2 != nil || shared2 != false || res2 != nil {
			t.Errorf("unexpected result from second call: %v, %v, %v", res2, err2, shared2)
		}
	})

	// 测试2: 请求失败，不共享结果
	t.Run("Test_Failed_Request_Not_Shared", func(t *testing.T) {
		var res any
		var err error
		var shared bool
		// 失败的请求
		res, err, shared = sf.Do("key2", func() (any, error) {
			fmt.Println("Running error function...")
			time.Sleep(1 * time.Second)
			return nil, fmt.Errorf("some error")
		})

		if err == nil || shared != false {
			t.Errorf("expected error and not shared, but got: %v, %v", res, shared)
		}
	})
}

func TestDoDupSuppress(t *testing.T) {
	var g sf
	var wg1, wg2 sync.WaitGroup
	c := make(chan string, 1)
	var calls int32
	fn := func() (interface{}, error) {
		if atomic.AddInt32(&calls, 1) == 1 {
			// First invocation.
			wg1.Done()
		}
		v := <-c
		c <- v // pump; make available for any future calls

		time.Sleep(10 * time.Millisecond) // let more goroutines enter Do

		return v, nil
	}

	const n = 10
	wg1.Add(1)
	for i := 0; i < n; i++ {
		wg1.Add(1)
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			wg1.Done()
			v, err, _ := g.Do("key", fn)
			if err != nil {
				t.Errorf("Do error: %v", err)
				return
			}
			if s, _ := v.(string); s != "bar" {
				t.Errorf("Do = %T %v; want %q", v, v, "bar")
			}
		}()
	}
	wg1.Wait()
	// At least one goroutine is in fn now and all of them have at
	// least reached the line before the Do.
	c <- "bar"
	wg2.Wait()
	if got := atomic.LoadInt32(&calls); got <= 0 || got >= n {
		t.Errorf("number of calls = %d; want over 0 and less than %d", got, n)
	}
}
