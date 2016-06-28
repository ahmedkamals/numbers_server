package app

import (
	"fmt"
	"runtime"
	"time"
)

type Dispatcher struct {
	workers []*Worker
	WorkersPool chan chan Job
	baseConfig map[string]string
	maxWorkers, maxQueuedItems, backendTimeout int
}

func NewDispatcher(maxWorkers, maxQueuedItems int, baseConfig map[string]string, backendTimeout int) *Dispatcher {

	pool := make(chan chan Job, maxWorkers)

	return &Dispatcher{
		workers: make([]*Worker, maxWorkers),
		WorkersPool: pool,
		maxWorkers: maxWorkers,
		maxQueuedItems: maxQueuedItems,
		baseConfig : baseConfig,
		backendTimeout: backendTimeout,
	}
}

func (self *Dispatcher) Run() {

	// Initializing the queue
	JobQueue = make(chan Job, self.maxQueuedItems)

	go self.stats()

	// Starting workers
	for i := 0; i < self.maxWorkers; i++ {
		self.workers[i] = NewWorker(i, self.WorkersPool)
		go self.workers[i].Start()
	}

	go self.dispatch()

	aggregator := NewAggregator()
	go aggregator.monitorNewData(self.backendTimeout)
	go aggregator.aggregate()

	// Launching number server
	numbersServer := NewNumbersServer()
	// Should be last line as it is a blocking.
	numbersServer.Start(self.baseConfig)
}

func (self *Dispatcher) dispatch() {

	for job := range JobQueue {

		go func(job Job) {

			// Blocking till an idle worker is available, try to obtain a worker job channel that is available.
			jobChannel := <-self.WorkersPool

			// Dispatch the job to the worker job channel.
			jobChannel <- job
		}(job)
	}
}

func (*Dispatcher) stats() {

	for {
		fmt.Println("Number of Go routines:", runtime.NumGoroutine())
		time.Sleep(10 * time.Second)
	}
}