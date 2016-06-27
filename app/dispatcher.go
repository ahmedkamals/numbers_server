package app

type Dispatcher struct {
	WorkersPool chan chan Job
	maxWorkers, maxQueuedItems  int
}

func NewDispatcher(maxWorkers, maxQueuedItems int) *Dispatcher {
	pool := make(chan chan Job, maxWorkers)
	return &Dispatcher{WorkersPool: pool, maxWorkers: maxWorkers, maxQueuedItems: maxQueuedItems}
}

func (self *Dispatcher) Run() {

	// Initializing the queue
	JobQueue = make(chan Job, self.maxQueuedItems)

	for i := 0; i < self.maxWorkers; i++ {
		worker := NewWorker(i, self.WorkersPool)
		worker.Start()
	}

	go self.dispatch()
}

func (self *Dispatcher) dispatch() {

	for {
		select {
		case job := <- JobQueue:
			go func(job Job) {
				// Blocking till an idle worker is available, try to obtain a worker job channel that is available.
				jobChannel := <-self.WorkersPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}