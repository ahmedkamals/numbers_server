package app

import (
	"../communication"
_	"../config"
	"time"
	"fmt"
)

type Worker struct {
	id         int
	jobRequest chan Job
	jobsPool   chan chan Job
	quit       chan bool
}

func NewWorker(id int, jobsPool chan chan Job) *Worker {
	return &Worker{
		id: id,
		jobRequest: make(chan Job),
		jobsPool: jobsPool,
		quit: make(chan bool),
	}
}

func (self *Worker) Start() {

	// Todo: use service locator
	//serviceLocator := &config.ServiceLocator{}
	//logger := serviceLocator.Logger()
	//
	//logger.Info("Worker started", logger)
	fmt.Println("Worker ", self.id, " started")

	for {
		// Register the current worker into the worker queue.
		self.jobsPool <- self.jobRequest

		select {
		// A work request is received.
		case job := <- self.jobRequest:
			self.process(job)

		// Workers will stop working after 24 hours, taking a nap :P
		case <-time.After(time.Hour * 24):
			self.Stop()
		case <-self.quit:
			return
		}
	}

}

func (self *Worker) Stop() {
	go func() {
		fmt.Println("Worker stopped ", self.id)
		self.quit <- true
	}()
}

func (self *Worker) process(job Job) {

	request := communication.NewRequest("", map[string]string{
		"method": job.Payload.method,
		"host": job.Payload.host,
		"path": job.Payload.path,
	})

	response, err := job.Payload.Fetch(request)

	fmt.Println("Worker ", self.id, "is processing Job ", job.Id())

	if nil != err {
		// Todo: using logger
		//logger.Error(err.Error(), nil)
		fmt.Println(err.Error())
		return
	}

	responsePayloadProcessor := ResponsePayloadProcessor{}
	numbers, err := responsePayloadProcessor.Process(response)

	if nil != err {
		// Todo: using logger
		fmt.Println(err.Error())
		return
	}

	// Pushing numbers to merge channel.
	mergeQueue <- numbers
}