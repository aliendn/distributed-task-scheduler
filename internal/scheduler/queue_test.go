package scheduler

import (
	"testing"
	"time"
)

func TestPriorityOrder(t *testing.T) {
	q := NewPriorityQueue()

	t1 := &Task{ID: "low", Priority: Low, CreatedAt: time.Now()}
	t2 := &Task{ID: "high", Priority: High, CreatedAt: time.Now().Add(-1 * time.Minute)}
	t3 := &Task{ID: "medium", Priority: Medium, CreatedAt: time.Now()}

	q.PushTask(t1)
	q.PushTask(t2)
	q.PushTask(t3)

	first := q.PopTask()
	if first.ID != "high" {
		t.Fatalf("expected high priority task, got %s", first.ID)
	}
	second := q.PopTask()
	if second.ID != "medium" {
		t.Fatalf("expected medium priority task, got %s", second.ID)
	}
	third := q.PopTask()
	if third.ID != "low" {
		t.Fatalf("expected low priority task, got %s", third.ID)
	}
}
