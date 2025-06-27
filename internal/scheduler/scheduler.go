package scheduler

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"distributed-task-scheduler/pkg/models"
	"distributed-task-scheduler/pkg/repositories"

	"github.com/google/uuid"
)

// TaskScheduler coordinates the queue + DB repo
type TaskScheduler struct {
	queue *PriorityQueue
	repo  *repositories.TaskRepository

	cache      map[string]*Task
	cacheMutex sync.RWMutex
}

// NewTaskScheduler binds queue + repo
func NewTaskScheduler(queue *PriorityQueue, repo *repositories.TaskRepository) *TaskScheduler {
	return &TaskScheduler{
		queue: queue,
		repo:  repo,
		cache: make(map[string]*Task),
	}
}

// SubmitTask persists + enqueues
func (ts *TaskScheduler) SubmitTask(priority TaskPriority, payload json.RawMessage) *Task {
	// Create Task
	task := &Task{
		ID:        uuid.New().String(),
		Priority:  priority,
		Payload:   payload,
		CreatedAt: time.Now().UTC(),
		Status:    "pending",
	}

	// Persist to DB
	dbTask := &models.Task{
		ID:        task.ID,
		Priority:  models.TaskPriority(priority),
		Payload:   payload,
		CreatedAt: task.CreatedAt,
		Status:    task.Status,
	}

	if err := ts.repo.Create(dbTask); err != nil {
		log.Printf("[Scheduler] DB insert failed: %v", err)
	}

	// Save to cache
	ts.cacheMutex.Lock()
	ts.cache[task.ID] = task
	ts.cacheMutex.Unlock()

	// Enqueue
	ts.queue.PushTask(task)

	log.Printf("[Scheduler] Submitted task %s with %s priority", task.ID, priority.String())
	return task
}

// GetTask gets from cache or DB fallback
func (ts *TaskScheduler) GetTask(id string) (*Task, bool) {
	ts.cacheMutex.RLock()
	if task, ok := ts.cache[id]; ok {
		ts.cacheMutex.RUnlock()
		return task, true
	}
	ts.cacheMutex.RUnlock()

	// Fallback: try DB
	dbTask, err := ts.repo.GetByID(id)
	if err == nil && dbTask != nil {
		task := &Task{
			ID:        dbTask.ID,
			Priority:  TaskPriority(dbTask.Priority),
			Payload:   dbTask.Payload,
			CreatedAt: dbTask.CreatedAt,
			Status:    dbTask.Status,
		}

		// Re-cache
		ts.cacheMutex.Lock()
		ts.cache[id] = task
		ts.cacheMutex.Unlock()

		return task, true
	}

	return nil, false
}

// RecoverUnfinishedTasks reloads from DB on startup
func (ts *TaskScheduler) RecoverUnfinishedTasks() {
	tasks, err := ts.repo.GetUnfinishedTasks()
	if err != nil {
		log.Printf("[Scheduler] Failed recovery query: %v", err)
		return
	}

	for _, dbTask := range tasks {
		task := &Task{
			ID:        dbTask.ID,
			Priority:  TaskPriority(dbTask.Priority),
			Payload:   dbTask.Payload,
			CreatedAt: dbTask.CreatedAt,
			Status:    dbTask.Status,
		}

		ts.cacheMutex.Lock()
		ts.cache[task.ID] = task
		ts.cacheMutex.Unlock()

		ts.queue.PushTask(task)
	}

	log.Printf("[Scheduler] Recovered %d unfinished tasks", len(tasks))
}
