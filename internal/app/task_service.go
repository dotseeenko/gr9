package app

import "github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"

type TaskService interface {
}

type taskService struct {
	taskRepo database.TaskRepository
}

func NewTaskService(tr database.TaskRepository) TaskService {
	return taskService{
		taskRepo: tr,
	}
}
