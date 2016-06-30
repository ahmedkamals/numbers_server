package processing

import (
	"time"
	"../services"
)

var (
	// A buffered channel that we can merge fetched data on.
	MergeQueue chan []int
	// A boolean chanel to indicate that merge is done, and the aggregation process should start.
	isMergeDone chan bool
	// A channel to serve aggregated data.
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

		// Giving extra 100ms for processing.
		case <-time.After(time.Millisecond * time.Duration(timeout - 100)):
			serviceLocator := services.NewServiceLocator()
			serviceLocator.Logger().Info("Timed out:", aggregatedData)

			isMergeDone <- true
		}
	}
}

func (self *Aggregator) Aggregate(operator *Operator) {

	for isMergeDone := range isMergeDone {

		if (isMergeDone) {

			processedData := operator.Process(aggregatedData)

			AggregationQueue <- processedData
		}
	}
}