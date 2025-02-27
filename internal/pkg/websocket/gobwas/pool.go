package gobwas

import (
	"sync"
	"time"
)

// Pool implements a worker pool pattern for goroutine reuse.
// It manages a fixed number of worker goroutines that process tasks from a queue.
type Pool struct {
	maxWorkers int        // Maximum number of workers allowed
	queue      chan task  // Queue for pending tasks
	semaphore  chan token // Semaphore to control worker count
	waitGroup  sync.WaitGroup
}

// Simplified type aliases for better readability
type (
	task  = func()
	token = struct{}
)

// NewPool creates a new goroutine pool with specified parameters.
//
// Parameters:
//   - maxWorkers: Maximum number of goroutines that can be created
//   - queueSize: Maximum number of tasks that can wait in the queue
//   - preAllocWorkers: Number of workers to create in advance (0 for lazy initialization)
//
// The pool will create workers on demand up to maxWorkers.
func NewPool(maxWorkers, queueSize, preAllocWorkers int) *Pool {
	pool := &Pool{
		maxWorkers: maxWorkers,
		queue:      make(chan task, queueSize),
		semaphore:  make(chan token, maxWorkers),
	}

	// Initialize workers upfront if requested
	if preAllocWorkers > 0 {
		preAllocWorkers = min(preAllocWorkers, maxWorkers)
		for i := 0; i < preAllocWorkers; i++ {
			pool.startWorker()
		}
	}

	return pool
}

// Schedule adds a task to be executed by the worker pool.
// If the queue is full, it tries to start a new worker.
// If the worker limit is reached, it blocks until a worker becomes available.
func (p *Pool) Schedule(task task) {
	// Try to add task to queue without blocking
	select {
	case p.queue <- task:
		return // Task scheduled successfully
	default:
		// Queue is full, try to spawn a new worker
	}

	// Try to acquire worker slot
	select {
	case p.semaphore <- token{}:
		// Worker slot acquired, start a new worker and queue the task
		p.startWorker()
		p.queue <- task
	default:
		// Worker limit reached, block until queue has space
		p.queue <- task
	}
}

// ScheduleTimeout attempts to schedule a task with a timeout.
// Returns ErrScheduleTimeout if the task couldn't be scheduled within the given timeout.
func (p *Pool) ScheduleTimeout(timeout time.Duration, task task) error {
	// Fast path: try to enqueue directly
	select {
	case p.queue <- task:
		return nil // Task scheduled successfully
	default:
		// Queue is full, try other approaches
	}

	// Try to start a new worker
	select {
	case p.semaphore <- token{}:
		// Worker slot acquired, start a new worker and queue the task
		p.startWorker()
		p.queue <- task
		return nil
	default:
		// Worker limit reached, try to queue with timeout
		select {
		case p.queue <- task:
			return nil
		case <-time.After(timeout):
			return ErrScheduleTimeout
		}
	}
}

// startWorker launches a new worker goroutine that processes tasks from the queue.
// The worker automatically exits when the queue is closed or when there are no more tasks.
func (p *Pool) startWorker() {
	p.waitGroup.Add(1)

	go func() {
		defer p.waitGroup.Done()
		defer func() { <-p.semaphore }() // Release worker slot when done

		// Process tasks until queue is closed
		for task := range p.queue {
			task()
		}
	}()
}

// Close gracefully shuts down the pool by closing the queue and waiting for all workers to finish.
// No new tasks should be scheduled after calling Close.
func (p *Pool) Close() {
	close(p.queue)
	p.waitGroup.Wait()
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
