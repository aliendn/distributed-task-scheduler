package main

import (
	"distributed-task-scheduler/internal/cluster"
	"distributed-task-scheduler/internal/metrics"
	"distributed-task-scheduler/internal/routes"
	"distributed-task-scheduler/internal/scheduler"
	"distributed-task-scheduler/pkg/database"
	"distributed-task-scheduler/pkg/repositories"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

// @title Distributed Task Scheduler API
// @version 1.0
// @description API for scheduling tasks
// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	metrics.Init()

	// Init PostgreSQL with GORM
	database.InitGorm()
	db := database.DB

	// Init repository
	taskRepo := repositories.NewTaskRepository(db)

	// Init scheduler
	queue := scheduler.NewPriorityQueue()
	taskScheduler := scheduler.NewTaskScheduler(queue, taskRepo)

	// Init worker pool with repo too
	workerPool := scheduler.NewWorkerPool(queue, taskRepo, 4)

	// Recover tasks from DB
	taskScheduler.RecoverUnfinishedTasks()

	// Start workers
	workerPool.Start()
	defer workerPool.Stop()

	// Cluster logic
	leader := cluster.NewLeaderElector(func() {
		log.Println("[Cluster] I am the leader. I can assign tasks.")
	})
	leader.Start()
	defer leader.Stop()

	heartBeater := cluster.NewHeartbeater(leader.NodeID, 5*time.Second)
	heartBeater.Start()
	defer heartBeater.Stop()
	router := gin.Default()
	routes.RegisterRoutes(router, taskScheduler)

	log.Println("🚀 Server running at http://localhost:8080")
	router.Run(":8080")
}
