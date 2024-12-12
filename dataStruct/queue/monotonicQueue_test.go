package queue

import "testing"

func TestMonotonicQueue(t *testing.T) {
	nums := []int{1, 3, -1, -3, 5, 3, 6, 7}
	k := 3
	mq := MonotonicQueue{}
	res := []int{}
	for i := 0; i < len(nums); i++ {
		if i < k-1 {
			mq.Push(nums[i])
		} else {
			mq.Push(nums[i])
			res = append(res, mq.queue[0])
			mq.Pop(nums[i-k+1])
		}
	}
	if res[0] != 3 || res[1] != 3 || res[2] != 5 || res[3] != 5 || res[4] != 6 || res[5] != 7 {
		t.Errorf("MonotonicQueue() failed")
	}
}
