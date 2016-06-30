package app

import (
	"runtime"
	"time"
	"../queue"
	"../services"
)

type Dispatcher struct {
	workers                    []*queue.Worker
	WorkersPool                chan chan queue.Job
	config                     map[string]string
	maxWorkers, maxQueuedItems int
}

func NewDispatcher(maxWorkers, maxQueuedItems int, config map[string]string) *Dispatcher {

	pool := make(chan chan queue.Job, maxWorkers)

	return &Dispatcher{
		workers: make([]*queue.Worker, maxWorkers),
		WorkersPool: pool,
		maxWorkers: maxWorkers,
		maxQueuedItems: maxQueuedItems,
		config : config,
	}
}

func (self *Dispatcher) Run() {

	// Initializing the queue
	queue.JobsQueue = make(chan queue.Job, self.maxQueuedItems)

	go self.stats()

	// Starting workers
	for i := 0; i < self.maxWorkers; i++ {
		self.workers[i] = queue.NewWorker(i, self.WorkersPool)
		go self.workers[i].Start()
	}

	go self.dispatch()

	// Launching number server
	numbersServer := NewNumbersServer(self.config)

	// Should be last line as it is a blocking.
	numbersServer.Start()
}

func (self *Dispatcher) dispatch() {

	for job := range queue.JobsQueue {

		go func(job queue.Job) {

			// Blocking till an idle worker is available, try to obtain a worker job channel that is available.
			jobChannel := <-self.WorkersPool

			// Dispatch the job to the worker job channel.
			jobChannel <- job
		}(job)
	}
}

func (*Dispatcher) stats() {

	serviceLocator := services.ServiceLocator{}

	for {
		serviceLocator.Logger().Info("Number of Go routines:", runtime.NumGoroutine())
		time.Sleep(10 * time.Second)
	}
}