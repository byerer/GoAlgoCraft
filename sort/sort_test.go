package sort

import (
	"testing"
)

func TestQuickSort(t *testing.T) {
	nums := []int{3, 2, 1, 5, 4}
	QuickSort(nums, 0, len(nums)-1)
	for i := 0; i < len(nums); i++ {
		if nums[i] != i+1 {
			t.Errorf("QuickSort() failed")
		}
	}
}
