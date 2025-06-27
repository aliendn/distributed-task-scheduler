package scheduler

import (
	"encoding/json"
	"testing"
)

func TestPriorityQueueOrdering(t *testing.T) {
	q := NewPriorityQueue()

	task1 := NewTask(Low, json.RawMessage(`{"data":1}`))
	task2 := NewTask(High, json.RawMessage(`{"data":2}`))
	task3 := NewTask(Medium, json.RawMessage(`{"data":3}`))

	q.PushTask(task1)
	q.PushTask(task2)
	q.PushTask(task3)

	if q.Len() != 3 {
		t.Fatalf("Expected 3 tasks, got %d", q.Len())
	}

	t1 := q.PopTask()
	if t1.Priority != High {
		t.Fatalf("Expected high priority first, got %s", t1.Priority.String())
	}

	t2 := q.PopTask()
	if t2.Priority != Medium {
		t.Fatalf("Expected medium priority second, got %s", t2.Priority.String())
	}

	t3 := q.PopTask()
	if t3.Priority != Low {
		t.Fatalf("Expected low priority third, got %s", t3.Priority.String())
	}
}
