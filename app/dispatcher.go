package app

import (
	"fmt"
	"runtime"
	"time"
	"../queue"
	"../processing"
)

type Dispatcher struct {
	workers []*queue.Worker
	WorkersPool chan chan queue.Job
	baseConfig map[string]string
	maxWorkers, maxQueuedItems, backendTimeout int
}

func NewDispatcher(maxWorkers, maxQueuedItems int, baseConfig map[string]string, backendTimeout int) *Dispatcher {

	pool := make(chan chan queue.Job, maxWorkers)

	return &Dispatcher{
		workers: make([]*queue.Worker, maxWorkers),
		WorkersPool: pool,
		maxWorkers: maxWorkers,
		maxQueuedItems: maxQueuedItems,
		baseConfig : baseConfig,
		backendTimeout: backendTimeout,
	}
}

func (self *Dispatcher) Run() {

	// Initializing the queue
	queue.JobQueue = make(chan queue.Job, self.maxQueuedItems)

	go self.stats()

	// Starting workers
	for i := 0; i < self.maxWorkers; i++ {
		self.workers[i] = queue.NewWorker(i, self.WorkersPool)
		go self.workers[i].Start()
	}

	go self.dispatch()

	aggregator := processing.NewAggregator()
	go aggregator.MonitorNewData(self.backendTimeout)
	go aggregator.Aggregate()

	// Launching number server
	numbersServer := NewNumbersServer()
	// Should be last line as it is a blocking.
	numbersServer.Start(self.baseConfig)
}

func (self *Dispatcher) dispatch() {

	for job := range queue.JobQueue {

		go func(job queue.Job) {

			// Blocking till an idle worker is available, try to obtain a worker job channel that is available.
			jobChannel := <-self.WorkersPool

			// Dispatch the job to the worker job channel.
			jobChannel <- job
		}(job)
	}
}

func (*Dispatcher) stats() {

	for {
		// Todo: Use logger
		fmt.Println("Number of Go routines:", runtime.NumGoroutine())
		time.Sleep(10 * time.Second)
	}
}