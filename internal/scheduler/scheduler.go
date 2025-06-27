package scheduler

import (
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
func (ts *TaskScheduler) SubmitTask(priority TaskPriority, payload interface{}) *Task {
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

// GetAllTasks returns all tasks from cache and DB fallback.
func (ts *TaskScheduler) GetAllTasks() []*Task {
	// get all from cache
	ts.cacheMutex.RLock()
	cacheTasks := make(map[string]*Task, len(ts.cache))
	for id, t := range ts.cache {
		cacheTasks[id] = t
	}
	ts.cacheMutex.RUnlock()

	// get all tasks from DB
	dbTasks, err := ts.repo.GetAll()
	if err != nil {
		log.Printf("[Scheduler] Failed to get tasks from DB: %v", err)
	}

	// Combine: avoid duplicates by checking IDs
	allTasks := make([]*Task, 0, len(cacheTasks)+len(dbTasks))
	for _, t := range cacheTasks {
		allTasks = append(allTasks, t)
	}

	for _, dbTask := range dbTasks {
		if _, exists := cacheTasks[dbTask.ID]; !exists {
			t := &Task{
				ID:        dbTask.ID,
				Priority:  TaskPriority(dbTask.Priority),
				Payload:   dbTask.Payload,
				CreatedAt: dbTask.CreatedAt,
				Status:    dbTask.Status,
			}
			allTasks = append(allTasks, t)
		}
	}

	return allTasks
}
