package app

import (
	"fmt"
)

type Aggregator struct {
	data []int
}

func (self *Aggregator)Start(timeout int) {

	for {
		select {
		case items := <- AggregationQueue:
			fmt.Println(items)
			self.data = append(self.data, items...)
		}
	}
}

// A buffered channel that we can send response payloads on.
var AggregationQueue chan []int