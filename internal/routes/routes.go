package routes

import (
	"distributed-task-scheduler/internal/api"
	"distributed-task-scheduler/internal/scheduler"

	_ "distributed-task-scheduler/cmd/distributed-task-scheduler/docs" // swag docs

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(s *scheduler.TaskScheduler) *gin.Engine {
	router := gin.Default()
	h := api.NewAPIHandler(s)

	router.POST("/tasks", h.SubmitTask)
	router.GET("/tasks/:id", h.GetTask)

	// Prometheus
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
