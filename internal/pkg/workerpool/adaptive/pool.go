package workerpool

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// Pool implements an adaptive worker pool pattern for goroutine reuse.
// It dynamically scales between minimum and maximum worker counts based on load.
type Pool struct {
	maxWorkers    int        // Maximum number of workers allowed
	minWorkers    int        // Minimum number of workers to maintain
	queue         chan task  // Queue for pending tasks
	semaphore     chan token // Semaphore to control worker count
	waitGroup     sync.WaitGroup
	activeWorkers int32         // Current number of active workers (atomic)
	idleTimeout   time.Duration // How long workers stay idle before terminating
	mutex         sync.Mutex    // For coordinating worker count operations
}

// Simplified type aliases for better readability
type (
	task  = func()
	token = struct{}
)

// NewPoolWithAutoScale creates a new goroutine pool with auto-scaling capabilities.
//
// Parameters:
//   - maxWorkers: Maximum number of goroutines that can be created
//   - minWorkers: Minimum number of workers to maintain (even when idle)
//   - queueSize: Maximum number of tasks that can wait in the queue
//   - idleTimeout: How long workers stay idle before terminating (if above minWorkers)
//
// The pool will scale between minWorkers and maxWorkers based on load.
func NewPoolWithAutoScale(maxWorkers, minWorkers, queueSize int, idleTimeout time.Duration) *Pool {
	// Ensure minimum workers doesn't exceed maximum
	minWorkers = min(minWorkers, maxWorkers)

	pool := &Pool{
		maxWorkers:    maxWorkers,
		minWorkers:    minWorkers,
		queue:         make(chan task, queueSize),
		semaphore:     make(chan token, maxWorkers),
		idleTimeout:   idleTimeout,
		activeWorkers: 0,
	}

	// Initialize minimum workers upfront
	for i := 0; i < minWorkers; i++ {
		pool.startPermanentWorker()
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
		// Worker slot acquired, start a new adaptive worker and queue the task
		p.startAdaptiveWorker()
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
		p.startAdaptiveWorker()
		p.queue <- task
		return nil
	default:
		// Worker limit reached, try to queue with timeout
		select {
		case p.queue <- task:
			return nil
		case <-time.After(timeout):
			return errors.New("ErrScheduleTimeout")
		}
	}
}

// startPermanentWorker launches a worker that never terminates until pool closure.
// These form the minimum worker set that's always available.
func (p *Pool) startPermanentWorker() {
	p.waitGroup.Add(1)
	atomic.AddInt32(&p.activeWorkers, 1)

	go func() {
		defer p.waitGroup.Done()
		defer atomic.AddInt32(&p.activeWorkers, -1)
		defer func() { <-p.semaphore }() // Release worker slot when done

		// Acquire semaphore slot
		p.semaphore <- token{}

		// Process tasks indefinitely until queue is closed
		for task := range p.queue {
			task()
		}
	}()
}

// startAdaptiveWorker launches a worker that can terminate after being idle.
// These form the dynamic part of the worker pool that scales based on load.
func (p *Pool) startAdaptiveWorker() {
	p.waitGroup.Add(1)
	atomic.AddInt32(&p.activeWorkers, 1)

	go func() {
		defer p.waitGroup.Done()
		defer atomic.AddInt32(&p.activeWorkers, -1)
		defer func() { <-p.semaphore }() // Release worker slot when done

		for {
			// Wait for a task or timeout
			select {
			case task, ok := <-p.queue:
				if !ok {
					// Queue is closed, terminate worker
					return
				}

				// Execute the task
				task()

			case <-time.After(p.idleTimeout):
				// Check if we're above minimum worker count
				if p.canTerminateWorker() {
					return // Terminate this worker
				}
			}
		}
	}()
}

// canTerminateWorker determines if a worker can terminate based on minimum requirements.
// Uses lock to ensure accurate counting during rapid scale-down.
func (p *Pool) canTerminateWorker() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	currentWorkers := atomic.LoadInt32(&p.activeWorkers)

	// Only terminate if we're above minimum worker count
	if currentWorkers > int32(p.minWorkers) {
		return true
	}

	return false
}

// ActiveWorkerCount returns the current number of active workers
func (p *Pool) ActiveWorkerCount() int {
	return int(atomic.LoadInt32(&p.activeWorkers))
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
