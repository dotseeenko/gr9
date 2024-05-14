package domain

import "time"

type Task struct {
	Id          uint64
	UserId      uint64
	Title       string
	Description string
	Status      TaskStatus
	Deadline    *time.Time
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

type TaskStatus string

const (
	New        TaskStatus = "NEW"
	InProgress TaskStatus = "IN_PROGRESS"
	Done       TaskStatus = "DONE"
)
