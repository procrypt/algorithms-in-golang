/*
Make sure that you handle the queue size properly, i,e
size field in MovingAverage struct is the size allocated to your queue
you check if you have sufficient space in your queue before adding
any element to the data array.
If at some point len(data) == size, you remove the top element from the queue
and update the queue and the current sum by reoving the element from the queue and
substracting that value from the sum.
*/

package main

type MovingAverage struct {
	data []int
	sum  int
	size int
}

/** Initialize your data structure here. */
func Constructor(size int) MovingAverage {
	return MovingAverage{size: size}
}

func (this *MovingAverage) Next(val int) float64 {
	if len(this.data) < this.size {
		this.sum += val
		this.data = append(this.data, val)
	} else if len(this.data) == this.size {
		// substract the value of the this.data[0] from the current sum
		this.sum -= this.data[0]
		// remove the first element of the queue
		this.data = this.data[1:]
		// append new data to the queue
		this.data = append(this.data, val)
		// update the current sun
		this.sum += val
	}
	return float64(this.sum) / float64(len(this.data))
}
