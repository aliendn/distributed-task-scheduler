package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"distributed-task-scheduler/internal/routes"
	"distributed-task-scheduler/internal/scheduler"
	"github.com/gin-gonic/gin"
)

func TestSubmitAndQueryTask(t *testing.T) {
	gin.SetMode(gin.TestMode)

	queue := scheduler.NewPriorityQueue()
	s := scheduler.NewTaskScheduler(queue)

	router := routes.SetupRouter(s)

	// Submit a task
	taskBody := map[string]interface{}{
		"priority": "high",
		"payload":  map[string]string{"action": "send_email", "to": "user@example.com"},
	}
	jsonData, _ := json.Marshal(taskBody)

	req := httptest.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusAccepted {
		t.Fatalf("Expected 202 Accepted, got %d", resp.Code)
	}

	var created map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &created)
	if err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}
	id := created["id"].(string)

	// Give worker a moment to process it
	time.Sleep(1 * time.Second)

	// Query task
	req2 := httptest.NewRequest("GET", "/tasks/"+id, nil)
	resp2 := httptest.NewRecorder()

	router.ServeHTTP(resp2, req2)

	if resp2.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", resp2.Code)
	}

	var fetched map[string]interface{}
	_ = json.Unmarshal(resp2.Body.Bytes(), &fetched)
	if fetched["id"] != id {
		t.Fatalf("Expected task ID %s, got %v", id, fetched["id"])
	}
}
