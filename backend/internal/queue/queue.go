package queue

import (
	"context"
	"sync"
	"time"
)

// Job represents a job in the queue
type Job struct {
	ID        string
	Type      string
	Payload   interface{}
	Status    string
	CreatedAt time.Time
	StartedAt *time.Time
	EndedAt   *time.Time
	Result    interface{}
	Error     error
}

// Queue manages a queue of jobs
type Queue struct {
	jobs     map[string]*Job
	mu       sync.RWMutex
	maxSize  int
	timeout  time.Duration
	workers  int
	jobChan  chan *Job
	quit     chan struct{}
}

// New creates a new job queue
func New(maxSize, workers int, timeout time.Duration) *Queue {
	return &Queue{
		jobs:     make(map[string]*Job),
		maxSize:  maxSize,
		timeout:  timeout,
		workers:  workers,
		jobChan:  make(chan *Job, maxSize),
		quit:     make(chan struct{}),
	}
}

// Start starts the queue workers
func (q *Queue) Start(handler func(context.Context, *Job) error) {
	for i := 0; i < q.workers; i++ {
		go q.worker(i, handler)
	}
}

// worker processes jobs from the queue
func (q *Queue) worker(id int, handler func(context.Context, *Job) error) {
	for {
		select {
		case job := <-q.jobChan:
			q.processJob(job, handler)
		case <-q.quit:
			return
		}
	}
}

// processJob processes a single job
func (q *Queue) processJob(job *Job, handler func(context.Context, *Job) error) {
	ctx, cancel := context.WithTimeout(context.Background(), q.timeout)
	defer cancel()

	now := time.Now()
	job.StartedAt = &now
	job.Status = "running"

	if err := handler(ctx, job); err != nil {
		job.Error = err
		job.Status = "failed"
	} else {
		job.Status = "completed"
	}

	ended := time.Now()
	job.EndedAt = &ended
}

// Enqueue adds a job to the queue
func (q *Queue) Enqueue(job *Job) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.jobs) >= q.maxSize {
		return ErrQueueFull
	}

	job.CreatedAt = time.Now()
	job.Status = "pending"
	q.jobs[job.ID] = job

	select {
	case q.jobChan <- job:
		return nil
	default:
		return ErrQueueFull
	}
}

// GetJob retrieves a job by ID
func (q *Queue) GetJob(id string) (*Job, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	job, ok := q.jobs[id]
	return job, ok
}

// CancelJob cancels a pending job
func (q *Queue) CancelJob(id string) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	job, ok := q.jobs[id]
	if !ok {
		return ErrJobNotFound
	}

	if job.Status != "pending" {
		return ErrJobNotCancellable
	}

	job.Status = "cancelled"
	return nil
}

// GetStats returns queue statistics
func (q *Queue) GetStats() map[string]int {
	q.mu.RLock()
	defer q.mu.RUnlock()

	stats := map[string]int{
		"total":   len(q.jobs),
		"pending": 0,
		"running": 0,
		"done":    0,
	}

	for _, job := range q.jobs {
		switch job.Status {
		case "pending":
			stats["pending"]++
		case "running":
			stats["running"]++
		case "completed", "failed":
			stats["done"]++
		}
	}

	return stats
}

// Stop stops the queue
func (q *Queue) Stop() {
	close(q.quit)
}

// Errors
var (
	ErrQueueFull        = &QueueError{"queue is full"}
	ErrJobNotFound      = &QueueError{"job not found"}
	ErrJobNotCancellable = &QueueError{"job cannot be cancelled"}
)

// QueueError represents a queue error
type QueueError struct {
	Message string
}

func (e *QueueError) Error() string {
	return e.Message
}