package app

import "fmt"

type Dispatcher struct {
	WorkersPool chan chan Job
	baseConfig map[string]string
	maxWorkers, maxQueuedItems, backendTimeout int
}

func NewDispatcher(maxWorkers, maxQueuedItems int, baseConfig map[string]string, backendTimeout int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{
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
	MergeQueue = make(chan []int, self.maxQueuedItems)
	AggregatedData = make(chan []int)
	MergedData = []int{}
	IsMergeDone = make(chan bool)

	// Starting workers
	for i := 0; i < self.maxWorkers; i++ {
		worker := NewWorker(i, self.WorkersPool)
		worker.Start()
	}

	go self.dispatch()
	go self.aggregate()

	// Launching number server
	numbersServer := NewNumbersServer()
	numbersServer.startMergeChanel(self.backendTimeout)
	// Should be last line as it is a blocking.
	numbersServer.Start(self.baseConfig)
}

func (self *Dispatcher) dispatch() {

	for {
		select {
		case job := <- JobQueue:
			go func(job Job) {
				// Blocking till an idle worker is available, try to obtain a worker job channel that is available.
				jobChannel := <-self.WorkersPool

				// Dispatch the job to the worker job channel.
				jobChannel <- job
			}(job)
		}
	}
}

func (*Dispatcher) aggregate() {
	for {
		select {
		case isMergeDone := <- IsMergeDone:
			if (isMergeDone) {

				aggregator := Aggregator{}
				aggregatedData := aggregator.process(MergedData)
				fmt.Println("dispatcher agg:", aggregatedData)
				AggregatedData <- aggregatedData
			}
		}
	}
}