package sort

func QuickSort(nums []int, low, high int) {
	if low < high {
		pivot := paritition(nums, low, high)
		QuickSort(nums, low, pivot-1)
		QuickSort(nums, pivot+1, high)
	}
}

func paritition(nums []int, low, high int) int {
	pivot := nums[low]
	for low < high {
		for low < high && nums[high] > pivot {
			high--
		}
		nums[low] = nums[high]
		for low < high && nums[low] < pivot {
			low++
		}
		nums[high] = nums[low]
	}
	nums[low] = pivot
	return low
}
