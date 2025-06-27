package models

import (
	"time"
)

type TaskPriority int

const (
	High TaskPriority = iota
	Medium
	Low
)

type Task struct {
	ID        string       `gorm:"primaryKey" json:"id"`
	Priority  TaskPriority `gorm:"index" json:"priority"`
	Payload   interface{}  `json:"payload" gorm:"type:jsonb"`
	CreatedAt time.Time    `json:"created_at"`
	Status    string       `json:"status"` // pending, running, completed, failed
}
