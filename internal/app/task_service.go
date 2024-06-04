package app

import (
	"log"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/BohdanBoriak/boilerplate-go-back/internal/infra/database"
)

type TaskService interface {
	Save(t domain.Task) (domain.Task, error)
	GetForUser(uId uint64) ([]domain.Task, error)
	GetByID(id uint64) (domain.Task, error)
	DeleteByID(id uint64) error
	UpdateStatus(id uint64, status domain.TaskStatus) error
}

type taskService struct {
	taskRepo database.TaskRepository
}

func NewTaskService(tr database.TaskRepository) TaskService {
	return taskService{
		taskRepo: tr,
	}
}

func (s taskService) Save(t domain.Task) (domain.Task, error) {
	task, err := s.taskRepo.Save(t)
	if err != nil {
		log.Printf("TaskService -> Save: %s", err)
		return domain.Task{}, err
	}
	return task, nil
}

func (s taskService) GetForUser(uId uint64) ([]domain.Task, error) {
	tasks, err := s.taskRepo.GetByUserId(uId)
	if err != nil {
		log.Printf("TaskService -> GetForUser: %s", err)
		return nil, err
	}
	return tasks, nil
}

// GetByID знаходить завдання за його ID через репозиторій та повертає його у вигляді об'єкта domain.Task.
// Якщо завдання не знайдено або сталася помилка, логуються повідомлення та повертається пустий об'єкт domain.Task і помилка.
func (s taskService) GetByID(id uint64) (domain.Task, error) {
	task, err := s.taskRepo.GetByID(id) // Викликаємо метод GetByID репозиторію для отримання завдання.
	if err != nil {
		log.Printf("TaskService -> GetByID: %s", err) // Логування помилки.
		return domain.Task{}, err                     // Повертаємо пустий об'єкт та помилку.
	}
	return task, nil // Повертаємо знайдене завдання та nil (як індикатор відсутності помилки).
}

// DeleteByID видаляє завдання за його ID через репозиторій.
// Якщо сталася помилка під час видалення, логуються повідомлення та повертається ця помилка.
func (s taskService) DeleteByID(id uint64) error {
	err := s.taskRepo.DeleteByID(id) // Викликаємо метод DeleteByID репозиторію для видалення завдання.
	if err != nil {
		log.Printf("TaskService -> DeleteByID: %s", err) // Логування помилки.
		return err                                       // Повертаємо помилку.
	}
	return nil // Якщо все пройшло успішно, повертаємо nil.
}

// UpdateStatus оновлює статус завдання за його ID через репозиторій.
// Якщо сталася помилка під час оновлення, логуються повідомлення та повертається ця помилка.
func (s taskService) UpdateStatus(id uint64, status domain.TaskStatus) error {
	err := s.taskRepo.UpdateStatus(id, status) // Викликаємо метод UpdateStatus репозиторію для оновлення статусу завдання.
	if err != nil {
		log.Printf("TaskService -> UpdateStatus: %s", err) // Логування помилки.
		return err                                         // Повертаємо помилку.
	}
	return nil // Якщо все пройшло успішно, повертаємо nil.
}
