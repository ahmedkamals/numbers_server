package app

import (
	"../communication"
_	"../config"
	"time"
	"fmt"
)

type Worker struct {
	id 	    int
	WorkersPool chan chan Job
	JobChannel  chan Job
	quit        chan bool
}

func NewWorker(id int, workersPool chan chan Job) *Worker {
	return &Worker{
		id: id,
		WorkersPool: workersPool,
		JobChannel: make(chan Job),
		quit: make(chan bool),
	}
}

func (self *Worker) Start() {
	go func() {
		// Todo: use service locator
		//serviceLocator := &config.ServiceLocator{}
		//logger := serviceLocator.Logger()
		//
		//logger.Info("Worker started", logger)
		self.WorkersPool <- self.JobChannel
		fmt.Println("Worker ", self.id, " started")
		select {
		// A work request is received.
		case job := <- self.JobChannel:
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
			MergeQueue <- numbers

		// Workers will stop working after 24 hours, taking a nap :P
		case <-time.After(time.Hour * 24):
			self.Stop()
		case <-self.quit:
			return
		}
	}()
}

func (self *Worker) Stop() {
	go func() {
		fmt.Println("Worker stopped ", self.id)
		self.quit <- true
	}()
}