package app

import (
	"../communication"
	"../communication/protocols/http"
)

// A buffered channel that we can send work requests on.
var JobQueue chan Job

type Payload struct {
	method   string
	protocol *http.Protocol
	host     string
	path     string
}

func NewPayload(method string, protocol *http.Protocol, host, path string) *Payload {
	return &Payload{method, protocol, host, path}
}

func (self * Payload) Fetch(req *communication.Request) (*communication.Response, error) {
	return self.protocol.Send(req)
}

type Job struct {
	id string
	// Making the payload it self separated from the job.
	Payload *Payload
}

func NewJob(id string, payload *Payload) *Job {
	return &Job{id, payload}
}

type JobCollection struct {
	jobs []Job
}

func NewJobCollection(jobs []Job) *JobCollection{
	return &JobCollection{jobs}
}

func (self *Job) Id() string {
	return self.id
}

func PushToChanel (jobCollection *JobCollection) {

	for _, work := range jobCollection.jobs {

		// Push the work to the queue.
		JobQueue <- work
	}
}
