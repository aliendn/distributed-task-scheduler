package main

import (
	"distributed-task-scheduler/internal/cluster"
	"distributed-task-scheduler/internal/metrics"
	"distributed-task-scheduler/internal/routes"
	"distributed-task-scheduler/internal/scheduler"
	"distributed-task-scheduler/pkg/database"
	"distributed-task-scheduler/pkg/repositories"
	"log"
	"time"
)

func main() {
	metrics.Init()

	// Init PostgreSQL with GORM
	database.InitGorm()
	db := database.DB

	// Init repository
	taskRepo := repositories.NewTaskRepository(db)

	// Init scheduler components
	queue := scheduler.NewPriorityQueue()
	taskScheduler := scheduler.NewTaskScheduler(queue, taskRepo)

	taskScheduler.RecoverUnfinishedTasks()
	workerPool := scheduler.NewWorkerPool(queue, 4)
	workerPool.Start()
	defer workerPool.Stop()

	// Cluster logic
	leader := cluster.NewLeaderElector(func() {
		log.Println("[Cluster] I am the leader. I can assign tasks.")
	})
	leader.Start()
	defer leader.Stop()

	heartbeater := cluster.NewHeartbeater(leader.NodeID, 5*time.Second)
	heartbeater.Start()
	defer heartbeater.Stop()

	// Start API
	router := routes.SetupRouter(taskScheduler)
	log.Println("🚀 Server running at http://localhost:8080")
	router.Run(":8080")
}
