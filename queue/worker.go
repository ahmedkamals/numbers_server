package queue

import (
	"time"
	"runtime"
	"../processing"
	"../communication"
	"../services"
)

// Possible worker stats
const (
	PAUSED = 0
	RUNNING = 1
	STOPPED = 2
)

type Worker struct {
	id         int
	jobRequest chan Job
	jobsPool   chan chan Job
	state      chan int
}

func NewWorker(id int, jobsPool chan chan Job) *Worker {
	return &Worker{
		id: id,
		jobRequest: make(chan Job),
		jobsPool: jobsPool,
		state: make(chan int),
	}
}

func (self *Worker) Start() {

	serviceLocator := services.NewServiceLocator()
	logger := serviceLocator.Logger()

	go self.setState(RUNNING)

	for {
		// Register the current worker into the worker queue.
		self.jobsPool <- self.jobRequest

		select {
		// A work request is received.
		case job := <-self.jobRequest:
			err := self.process(job)
			if nil != err {
				logger.Error(err.Error())
			}

		// Workers will stop working after 24 hours, taking a nap :P
		case <-time.After(time.Hour * 24):
			self.Stop()

		case state := <-self.state:

			switch state {
			case PAUSED:
				logger.Info("Worker", self.id, "is paused.")

			case RUNNING:
				logger.Info("Worker", self.id, "is started.")

			case STOPPED:
				logger.Info("Worker", self.id, "is stopped.")

			default:
				// We use runtime.Gosched() to prevent a deadlock in this case.
				// It will not be needed of work is performed here which yields
				// to the scheduler.
				runtime.Gosched()

				if PAUSED == state {
					break
				}
			}
		}
	}

}

func (self *Worker) Pause() {

	go self.setState(PAUSED)
}

func (self *Worker) Stop() {

	go self.setState(STOPPED)
}

func (self *Worker) setState(status int) {

	self.state <- status
}

func (self *Worker) process(job Job) error {

	request := communication.NewRequest("", map[string]string{
		"method": job.Payload.method,
		"host": job.Payload.host,
		"path": job.Payload.path,
	})

	response, err := job.Payload.Fetch(request)

	if nil != err {

		return err
	}

	serviceLocator := services.NewServiceLocator()
	serviceLocator.Logger().Info("Worker ", self.id, "is processing Job", job.Id())

	responsePayloadProcessor := processing.ResponsePayloadProcessor{}
	numbers, err := responsePayloadProcessor.Process(response)

	if nil != err {

		return err
	}

	// Pushing numbers to merge channel.
	processing.MergeQueue <- numbers

	return nil
}