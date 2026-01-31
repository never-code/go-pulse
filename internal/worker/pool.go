package worker

import "log"

// Pool manages workers and job queue
type Pool struct {
	WorkerCount int
	Jobs        chan Job
}

// NewPool creates a worker pool
func NewPool(workerCount int) *Pool {
	return &Pool{
		WorkerCount: workerCount,
		Jobs:        make(chan Job),
	}
}

// Start launches worker goroutines
func (p *Pool) Start() {
	for i := 1; i <= p.WorkerCount; i++ {
		go p.worker(i)
	}
}

// worker consumes jobs from the channel
func (p *Pool) worker(id int) {
	log.Printf("Worker %d started", id)

	for job := range p.Jobs {
		log.Printf("Worker %d processing URL: %s", id, job.URL)
	}
}

// Submit sends a job to the pool
func (p *Pool) Submit(job Job) {
	p.Jobs <- job
}
