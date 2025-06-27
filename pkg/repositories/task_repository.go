package repositories

import (
	"distributed-task-scheduler/pkg/models"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) UpdateStatus(id string, status string) error {
	return r.db.Model(&models.Task{}).Where("id = ?", id).Update("status", status).Error
}

func (r *TaskRepository) GetByID(id string) (*models.Task, error) {
	var task models.Task
	err := r.db.First(&task, "id = ?", id).Error
	return &task, err
}

func (r *TaskRepository) GetUnfinishedTasks() ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Where("status IN ?", []string{"pending", "running"}).Find(&tasks).Error
	return tasks, err
}

// GetAll returns all tasks in the database
func (r *TaskRepository) GetAll() ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}
