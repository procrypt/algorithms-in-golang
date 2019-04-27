package main

type MinStack struct {
	data []int
	min  []int
}

/** initialize your data structure here. */
func Constructor() MinStack {
	return MinStack{}
}

func (ms *MinStack) Push(x int) {
	ms.data = append(ms.data, x)
	if len(ms.min) == 0 {
		ms.min = append(ms.min, x)
	} else if ms.min[len(ms.min)-1] >= x {
		ms.min = append(ms.min, x)
	}
}

func (ms *MinStack) Pop() {
	if len(ms.data) > 0 {
		pop := ms.data[len(ms.data)-1]
		ms.data = ms.data[:len(ms.data)-1]
		if pop == ms.min[len(ms.min)-1] {
			ms.min = ms.min[:len(ms.min)-1]
		}
	}
}

func (ms *MinStack) Top() int {
	if len(ms.data) > 0 {
		return ms.data[len(ms.data)-1]
	}
	return 0
}

func (ms *MinStack) GetMin() int {
	if len(ms.min) > 0 {
		return ms.min[len(ms.min)-1]
	}
	return 0
}
