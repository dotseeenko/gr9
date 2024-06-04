package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/app"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/requests"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/http/resources"
	"github.com/go-chi/chi/v5"
)

type TaskController struct {
	taskService app.TaskService
}

func NewTaskController(ts app.TaskService) TaskController {
	return TaskController{
		taskService: ts,
	}
}

func (c TaskController) Save() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, err := requests.Bind(r, requests.TaskRequest{}, domain.Task{})
		if err != nil {
			log.Printf("TaskController -> Save: %s", err)
			BadRequest(w, err)
			return
		}

		user := r.Context().Value(UserKey).(domain.User)
		task.UserId = user.Id
		task.Status = domain.New
		task, err = c.taskService.Save(task)
		if err != nil {
			log.Printf("TaskController -> Save: %s", err)
			InternalServerError(w, err)
			return
		}

		var tDto resources.TaskDto
		tDto = tDto.DomainToDto(task)
		Created(w, tDto)
	}
}

func (c TaskController) GetForUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.Context().Value(UserKey).(domain.User)

		tasks, err := c.taskService.GetForUser(user.Id)
		if err != nil {
			log.Printf("TaskController -> GetForUser: %s", err)
			InternalServerError(w, err)
			return
		}

		var tasksDto resources.TasksDto
		tasksDto = tasksDto.DomainToDtoCollection(tasks)
		Success(w, tasksDto)
	}
}

// GetByID повертає HTTP-обробник, який обробляє запити для отримання завдання за його ID.
func (c TaskController) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Отримуємо ID з URL параметра та конвертуємо його в uint64.
		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			log.Printf("TaskController -> GetByID: %s", err)
			BadRequest(w, err) // Відправляємо клієнту відповідь про неправильний запит.
			return
		}

		// Викликаємо метод сервісу для отримання завдання за ID.
		task, err := c.taskService.GetByID(id)
		if err != nil {
			log.Printf("TaskController -> GetByID: %s", err)
			InternalServerError(w, err) // Відправляємо клієнту відповідь про внутрішню помилку сервера.
			return
		}

		// Перетворюємо доменну модель завдання в DTO для відправки клієнту.
		var tDto resources.TaskDto
		tDto = tDto.DomainToDto(task)
		Success(w, tDto) // Відправляємо успішну відповідь з даними завдання.
	}
}

// DeleteByID повертає HTTP-обробник, який обробляє запити для видалення завдання за його ID.
func (c TaskController) DeleteByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Отримуємо ID з URL параметра та конвертуємо його в uint64.
		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			log.Printf("TaskController -> DeleteByID: %s", err)
			BadRequest(w, err) // Відправляємо клієнту відповідь про неправильний запит.
			return
		}

		// Викликаємо метод сервісу для видалення завдання за ID.
		err = c.taskService.DeleteByID(id)
		if err != nil {
			log.Printf("TaskController -> DeleteByID: %s", err)
			InternalServerError(w, err) // Відправляємо клієнту відповідь про внутрішню помилку сервера.
			return
		}

		Success(w, "Завдання успішно видалено!") // Відправляємо успішну відповідь про видалення завдання.
	}
}

// UpdateStatus повертає HTTP-обробник, який обробляє запити для оновлення статусу завдання за його ID.
func (c TaskController) UpdateStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Отримуємо ID з URL параметра та конвертуємо його в uint64.
		id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
		if err != nil {
			log.Printf("TaskController -> UpdateStatus: %s", err)
			BadRequest(w, err) // Відправляємо клієнту відповідь про неправильний запит.
			return
		}

		// Структура для зчитування нового статусу із запиту.
		var statusUpdateReq struct {
			Status domain.TaskStatus `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&statusUpdateReq); err != nil {
			log.Printf("TaskController -> UpdateStatus: %s", err)
			BadRequest(w, err) // Відправляємо клієнту відповідь про неправильний запит.
			return
		}

		// Перевіряємо, чи є новий статус дійсним.
		if !isValidTaskStatus(statusUpdateReq.Status) {
			err := fmt.Errorf("Не дійсний статус: %s. Використайте будь-ласка: NEW or IN_PROGRESS or DONE", statusUpdateReq.Status)
			log.Printf("TaskController -> UpdateStatus: %s", err)
			BadRequest(w, err) // Відправляємо клієнту відповідь про неправильний запит
			return
		}

		// Викликаємо метод сервісу для оновлення статусу завдання за ID.
		err = c.taskService.UpdateStatus(id, statusUpdateReq.Status)
		if err != nil {
			log.Printf("TaskController -> UpdateStatus: %s", err)
			InternalServerError(w, err) // Відправляємо клієнту відповідь про внутрішню помилку сервера.
			return
		}

		Success(w, "Статус успішно змінено!") // Відправляємо успішну відповідь про оновлення статусу завдання.
	}
}

// isValidTaskStatus перевіряє, чи є наданий статус дійсним.
func isValidTaskStatus(status domain.TaskStatus) bool {
	switch status {
	case domain.New, domain.InProgress, domain.Done:
		return true
	default:
		return false
	}
}
