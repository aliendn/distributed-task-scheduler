package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"distributed-task-scheduler/internal/metrics"
	"distributed-task-scheduler/pkg/repositories"
)

// WorkerPool runs N workers.
type WorkerPool struct {
	queue     *PriorityQueue
	repo      *repositories.TaskRepository
	workerNum int
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewWorkerPool with repo for DB updates.
func NewWorkerPool(queue *PriorityQueue, repo *repositories.TaskRepository, workerNum int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		queue:     queue,
		repo:      repo,
		workerNum: workerNum,
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (wp *WorkerPool) Start() {
	for i := 0; i < wp.workerNum; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
	log.Printf("[WorkerPool] Started %d workers", wp.workerNum)
}

func (wp *WorkerPool) Stop() {
	log.Println("[WorkerPool] Stopping...")
	wp.cancel()
	wp.wg.Wait()
	log.Println("[WorkerPool] All workers stopped.")
}

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

func (wp *WorkerPool) processTask(workerID int, task *Task) {
	log.Printf("[Worker %d] Processing task %s (Priority: %s)", workerID, task.ID, task.Priority.String())

	start := time.Now()

	// Mark as running
	task.Status = "running"
	if err := wp.repo.UpdateStatus(task.ID, task.Status); err != nil {
		log.Printf("[Worker %d] Failed DB update: %v", workerID, err)
	}

	time.Sleep(2 * time.Second) // Simulated work

	// Mark as completed
	task.Status = "completed"
	if err := wp.repo.UpdateStatus(task.ID, task.Status); err != nil {
		log.Printf("[Worker %d] Failed DB update: %v", workerID, err)
	}

	duration := time.Since(start).Seconds()
	metrics.TaskDuration.WithLabelValues(task.Priority.String()).Observe(duration)
	metrics.TasksProcessed.WithLabelValues(task.Status).Inc()

	log.Printf("[Worker %d] Completed task %s in %.2fs", workerID, task.ID, duration)
}
