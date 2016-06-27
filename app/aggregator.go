package app

import (
	"sort"
)

// A channel to serve aggregated data
var AggregatedData chan []int

type Aggregator struct {
}

func (self *Aggregator) process(data []int) []int {

	if 0 == len(data) {
		return data
	}

	data = self.removeDuplicates(data)
	data = self.sort(data)
	return data
}

func (*Aggregator) removeDuplicates(data []int) []int {

	encountered := map[int]bool{}
	result := []int{}

	for _, value := range data {
		if !encountered[value]{
			encountered[value] = true
			result = append(result, value)
		}
	}

	return result
}

func (*Aggregator) sort(data []int) []int {
	sort.Ints(data)
	return data
}

