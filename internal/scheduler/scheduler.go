package scheduler

import (
	"encoding/json"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

// TaskScheduler manages task intake and coordination with the worker pool.
type TaskScheduler struct {
	queue      *PriorityQueue
	taskStore  map[string]*Task // in-memory store, replace with DB in production
	storeMutex sync.RWMutex
}

// NewTaskScheduler creates a scheduler instance.
func NewTaskScheduler(queue *PriorityQueue) *TaskScheduler {
	return &TaskScheduler{
		queue:     queue,
		taskStore: make(map[string]*Task),
	}
}

// SubmitTask creates, stores, and enqueues a new task.
func (ts *TaskScheduler) SubmitTask(priority TaskPriority, payload json.RawMessage) *Task {
	task := &Task{
		ID:        uuid.New().String(),
		Priority:  priority,
		Payload:   payload,
		CreatedAt: time.Now().UTC(),
		Status:    "pending",
	}

	ts.storeMutex.Lock()
	ts.taskStore[task.ID] = task
	ts.storeMutex.Unlock()

	ts.queue.PushTask(task)

	log.Printf("[Scheduler] Submitted task %s with %s priority", task.ID, priority.String())
	return task
}

// GetTask retrieves task status by ID.
func (ts *TaskScheduler) GetTask(id string) (*Task, bool) {
	ts.storeMutex.RLock()
	defer ts.storeMutex.RUnlock()

	task, exists := ts.taskStore[id]
	return task, exists
}

// RecoverUnfinishedTasks (stub): Call this after a restart to requeue uncompleted tasks.
// In production, would fetch from persistent store.
func (ts *TaskScheduler) RecoverUnfinishedTasks() {
	// Example logic: loop over ts.taskStore and requeue unfinished
	log.Println("[Scheduler] Recovery placeholder - implement persistence later")
}

func TestPriorityQueue_PushPop(t *testing.T) {
	q := NewPriorityQueue()

	task1 := &Task{ID: "1", Priority: High, CreatedAt: time.Now()}
	task2 := &Task{ID: "2", Priority: Medium, CreatedAt: time.Now()}

	q.PushTask(task2)
	q.PushTask(task1)

	if q.Len() != 2 {
		t.Fatalf("expected length 2, got %d", q.Len())
	}

	popped := q.PopTask()
	if popped.ID != "1" {
		t.Fatalf("expected high priority task first, got %s", popped.ID)
	}
}
