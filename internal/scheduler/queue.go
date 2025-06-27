package scheduler

import (
	"container/heap"
	"sync"
	"time"

	"distributed-task-scheduler/internal/metrics"
	"github.com/google/uuid"
)

// TaskPriority defines task priority levels
type TaskPriority int

const (
	High TaskPriority = iota
	Medium
	Low
)

func (p TaskPriority) String() string {
	switch p {
	case High:
		return "high"
	case Medium:
		return "medium"
	case Low:
		return "low"
	default:
		return "unknown"
	}
}

// Task represents a unit of work
type Task struct {
	ID        string       `json:"id"`
	Priority  TaskPriority `json:"priority"`
	Payload   interface{}  `json:"payload"`
	CreatedAt time.Time    `json:"created_at"`
	Status    string       `json:"status"`
}

// TaskQueueItem wraps a Task for use in a heap
type TaskQueueItem struct {
	Task     *Task
	index    int
	priority TaskPriority
	created  time.Time
}

// PriorityQueue is a threadsafe min-heap by priority
type PriorityQueue struct {
	items []*TaskQueueItem
	lock  sync.Mutex
	cond  *sync.Cond
}

func NewPriorityQueue() *PriorityQueue {
	pq := &PriorityQueue{
		items: make([]*TaskQueueItem, 0),
	}
	pq.cond = sync.NewCond(&pq.lock)
	heap.Init(pq)
	return pq
}

func (pq *PriorityQueue) Len() int {
	pq.lock.Lock()
	defer pq.lock.Unlock()
	return len(pq.items)
}

func (pq *PriorityQueue) PushTask(task *Task) {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	item := &TaskQueueItem{
		Task:     task,
		priority: task.Priority,
		created:  task.CreatedAt,
	}
	heap.Push(pq, item)
	metrics.TasksInQueue.Inc()
	pq.cond.Signal()
}

func (pq *PriorityQueue) PopTask() *Task {
	pq.lock.Lock()
	defer pq.lock.Unlock()

	for len(pq.items) == 0 {
		pq.cond.Wait()
	}

	item := heap.Pop(pq).(*TaskQueueItem)
	metrics.TasksInQueue.Dec()
	return item.Task
}

func (pq PriorityQueue) Less(i, j int) bool {
	if pq.items[i].priority == pq.items[j].priority {
		return pq.items[i].created.Before(pq.items[j].created)
	}
	return pq.items[i].priority < pq.items[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*TaskQueueItem)
	item.index = len(pq.items)
	pq.items = append(pq.items, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := pq.items
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	pq.items = old[:n-1]
	return item
}

func NewTask(priority TaskPriority, payload interface{}) *Task {
	return &Task{
		ID:        uuid.New().String(),
		Priority:  priority,
		Payload:   payload,
		CreatedAt: time.Now().UTC(),
		Status:    "pending",
	}
}
