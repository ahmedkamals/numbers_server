package queue

// A buffered channel that we can send work requests on.
var JobsQueue chan Job

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
		JobsQueue <- work
	}
}
