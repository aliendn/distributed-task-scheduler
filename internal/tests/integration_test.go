package tests

import (
	"bytes"
	"distributed-task-scheduler/pkg/database"
	"distributed-task-scheduler/pkg/repositories"
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
	db := database.DB

	taskRepo := repositories.NewTaskRepository(db)

	queue := scheduler.NewPriorityQueue()
	s := scheduler.NewTaskScheduler(queue, taskRepo)

	router := gin.New()
	routes.RegisterRoutes(router, s)

	// Submit a task - note the full API prefix /api/v1/tasks
	taskBody := map[string]interface{}{
		"priority": "high",
		"payload":  map[string]string{"action": "send_email", "to": "user@example.com"},
	}
	jsonData, _ := json.Marshal(taskBody)

	req := httptest.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(jsonData))
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
	id, ok := created["id"].(string)
	if !ok {
		t.Fatalf("Response missing task ID")
	}

	// Give worker a moment to process it (if needed)
	time.Sleep(1 * time.Second)

	// Query task - use correct path with /api/v1/tasks/{id}
	req2 := httptest.NewRequest("GET", "/api/v1/tasks/"+id, nil)
	resp2 := httptest.NewRecorder()

	router.ServeHTTP(resp2, req2)

	if resp2.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", resp2.Code)
	}

	var fetched map[string]interface{}
	err = json.Unmarshal(resp2.Body.Bytes(), &fetched)
	if err != nil {
		t.Fatalf("Failed to parse fetched task: %v", err)
	}
	if fetched["id"] != id {
		t.Fatalf("Expected task ID %s, got %v", id, fetched["id"])
	}
}
