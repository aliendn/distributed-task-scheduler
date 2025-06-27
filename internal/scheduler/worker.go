package scheduler

import (
	"context"
	"log"
	"sync"
	"time"
)

// WorkerPool represents a pool of workers that process tasks.
type WorkerPool struct {
	queue      *PriorityQueue
	workerNum  int
	wg         sync.WaitGroup
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// NewWorkerPool initializes the pool.
func NewWorkerPool(queue *PriorityQueue, workerNum int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		queue:      queue,
		workerNum:  workerNum,
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

// Start begins all workers.
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerNum; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	log.Printf("[WorkerPool] Started %d workers", wp.workerNum)
}

// Stop gracefully shuts down the pool.
func (wp *WorkerPool) Stop() {
	log.Println("[WorkerPool] Stopping workers...")
	wp.cancelFunc()
	wp.wg.Wait()
	log.Println("[WorkerPool] All workers stopped.")
}

// worker continuously fetches and processes tasks.
func (wp *WorkerPool) worker(id int) {
	defer wp.wg.Done()
	for {
		select {
		case <-wp.ctx.Done():
			log.Printf("[Worker %d] Shutting down", id)
			return
		default:
			task := wp.queue.PopTask()
			if task == nil {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			wp.processTask(id, task)
		}
	}
}

// processTask simulates task processing.
func (wp *WorkerPool) processTask(workerID int, task *Task) {
	log.Printf("[Worker %d] Processing task %s (Priority: %s)", workerID, task.ID, task.Priority.String())

	task.Status = "running"
	time.Sleep(time.Second * 2) // Simulate work

	// After task completes
	task.Status = "completed"
	log.Printf("[Worker %d] Completed task %s", workerID, task.ID)
}
