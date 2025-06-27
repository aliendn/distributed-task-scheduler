package routes

import (
	"distributed-task-scheduler/internal/api"
	"distributed-task-scheduler/internal/scheduler"

	_ "distributed-task-scheduler/docs"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes sets up all routes on the given router.
func RegisterRoutes(router *gin.Engine, s *scheduler.TaskScheduler) {
	h := api.NewAPIHandler(s)

	v1 := router.Group("/api/v1")
	{
		v1.POST("/tasks", h.SubmitTask)
		v1.GET("/tasks/:id", h.GetTask)
		v1.GET("/tasks", h.GetAllTasks) // ✅ new GET /tasks
	}

	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
