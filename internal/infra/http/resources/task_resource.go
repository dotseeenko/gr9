package resources

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
)

type TaskDto struct {
	Id          uint64            `json:"id"`
	UserId      uint64            `json:"userId"`
	Title       string            `json:"title"`
	Description string            `json:"description, omit"`
	Status      domain.TaskStatus `json:"status"`
	Deadline    *time.Time        `json:"deadline, omitempty"`
}

type TasksDto struct {
	Tasks []TaskDto `json: "tasks"`
}

func (d TaskDto) DomainToDto(t domain.Task) TaskDto {
	return TaskDto{
		Id:          t.Id,
		UserId:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Deadline:    t.Deadline,
	}
}

func (d TasksDto) DomainToDtoCollection(ts []domain.Task) TasksDto {
	var tasks []TaskDto
	for _, t := range ts {
		tDto := TaskDto{}.DomainToDto(t)
		tasks = append(tasks, tDto)
	}
	return TasksDto{Tasks: tasks}
}
