package processing

import (
	"sort"
	"time"
	"fmt"
)

var (
	// A buffered channel that we can merge fetched data on.
	MergeQueue chan []int
	isMergeDone chan bool
	// A channel to serve aggregated data
	AggregationQueue chan []int
	aggregatedData []int
)

type Aggregator struct {
}

func NewAggregator() *Aggregator {

	MergeQueue = make(chan []int)
	isMergeDone = make(chan bool)
	AggregationQueue = make(chan []int)
	aggregatedData = []int{}

	return &Aggregator{}
}

func (self *Aggregator) MonitorNewData(timeout int) {

	for {
		select {
		case items := <-MergeQueue:
			aggregatedData = append(aggregatedData, items...)

		// Giving extra 100ms for processing
		case <- time.After(time.Millisecond * time.Duration(timeout - 100)):
			fmt.Println("timed out", aggregatedData)
			isMergeDone <- true
		}
	}
}

func (self *Aggregator) Aggregate() {

	for isMergeDone := range isMergeDone {

		if (isMergeDone) {

			aggregatedData := self.Process(aggregatedData)

			AggregationQueue <- aggregatedData
		}
	}
}

func (self *Aggregator) Process(data []int) []int {

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

