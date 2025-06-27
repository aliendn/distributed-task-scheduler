package api

import (
	"encoding/json"
	"net/http"

	"distributed-task-scheduler/internal/scheduler"
	"github.com/gin-gonic/gin"
)

// TaskRequest represents the request payload for a new task
type TaskRequest struct {
	Priority string          `json:"priority" binding:"required" example:"high"`
	Payload  json.RawMessage `json:"payload" binding:"required" example:"{\"action\":\"send_email\"}"`
}

// APIHandler wraps dependencies like the scheduler
type APIHandler struct {
	Scheduler *scheduler.TaskScheduler
}

// NewAPIHandler returns an initialized handler
func NewAPIHandler(s *scheduler.TaskScheduler) *APIHandler {
	return &APIHandler{Scheduler: s}
}

// SubmitTask godoc
// @Summary Submit a new task
// @Description Submit a task with priority and JSON payload
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task body TaskRequest true "Task to submit"
// @Success 202 {object} scheduler.Task
// @Failure 400 {object} map[string]string
// @Router /api/v1/tasks [post]
func (h *APIHandler) SubmitTask(c *gin.Context) {
	var req TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var priority scheduler.TaskPriority
	switch req.Priority {
	case "high":
		priority = scheduler.High
	case "medium":
		priority = scheduler.Medium
	case "low":
		priority = scheduler.Low
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid priority (must be high, medium, or low)"})
		return
	}

	task := h.Scheduler.SubmitTask(priority, req.Payload)
	c.JSON(http.StatusAccepted, task)
}

// GetTask godoc
// @Summary Get task by ID
// @Description Returns task status
// @Tags Tasks
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} scheduler.Task
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [get]
func (h *APIHandler) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, exists := h.Scheduler.GetTask(id)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}
	c.JSON(http.StatusOK, task)
}
