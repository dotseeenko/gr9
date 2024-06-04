package database

import (
	"time"

	"github.com/BohdanBoriak/boilerplate-go-back/internal/domain"
	"github.com/upper/db/v4"
)

const TasksTableName = "tasks"

type task struct {
	Id          uint64            `db:"id,omitempty"`
	UserId      uint64            `db:"user_id"`
	Title       string            `db:"title"`
	Description string            `db:"description"`
	Status      domain.TaskStatus `db:"status"`
	Deadline    *time.Time        `db:"deadline"`
	CreatedDate time.Time         `db:"created_date"`
	UpdatedDate time.Time         `db:"updated_date"`
	DeletedDate *time.Time        `db:"deleted_date"`
}

type TaskRepository interface {
	Save(t domain.Task) (domain.Task, error)
	GetByUserId(uId uint64) ([]domain.Task, error)
	GetByID(id uint64) (domain.Task, error)
	DeleteByID(id uint64) error
	UpdateStatus(id uint64, status domain.TaskStatus) error
}

type taskRepository struct {
	coll db.Collection
	sess db.Session
}

func NewTaskRepository(session db.Session) TaskRepository {
	return taskRepository{
		coll: session.Collection(TasksTableName),
		sess: session,
	}
}

func (r taskRepository) Save(t domain.Task) (domain.Task, error) {
	tsk := r.mapDomainToModel(t)
	tsk.CreatedDate = time.Now()
	tsk.UpdatedDate = time.Now()
	err := r.coll.InsertReturning(&tsk)
	if err != nil {
		return domain.Task{}, err
	}
	t = r.mapModelToDomain(tsk)
	return t, nil
}

func (r taskRepository) GetByUserId(uId uint64) ([]domain.Task, error) {
	var tasks []task
	err := r.coll.
		Find(db.Cond{"user_id": uId, "deleted_date": nil}).
		All(&tasks)
	if err != nil {
		return nil, err
	}

	res := r.mapModelToDomainCollection(tasks)
	return res, nil
}

func (r taskRepository) mapDomainToModel(t domain.Task) task {
	return task{
		Id:          t.Id,
		UserId:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Deadline:    t.Deadline,
		CreatedDate: t.CreatedDate,
		UpdatedDate: t.UpdatedDate,
		DeletedDate: t.DeletedDate,
	}
}

func (r taskRepository) mapModelToDomain(t task) domain.Task {
	return domain.Task{
		Id:          t.Id,
		UserId:      t.UserId,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		Deadline:    t.Deadline,
		CreatedDate: t.CreatedDate,
		UpdatedDate: t.UpdatedDate,
		DeletedDate: t.DeletedDate,
	}
}

func (r taskRepository) mapModelToDomainCollection(ts []task) []domain.Task {
	var tasks []domain.Task
	for _, t := range ts {
		tasks = append(tasks, r.mapModelToDomain(t))
	}
	return tasks
}

// GetByID знаходить завдання за його ID та повертає його у вигляді об'єкта domain.Task.
// Якщо завдання не знайдено або сталася помилка, повертається пустий об'єкт domain.Task і помилка.
func (r taskRepository) GetByID(id uint64) (domain.Task, error) {
	var tsk task                                                         // Створюємо змінну для зберігання знайденого завдання.
	err := r.coll.Find(db.Cond{"id": id, "deleted_date": nil}).One(&tsk) // Виконуємо пошук завдання за ID і відсутності дати видалення.
	if err != nil {
		return domain.Task{}, err // Якщо сталася помилка, повертаємо пустий об'єкт і помилку.
	}
	return r.mapModelToDomain(tsk), nil // Перетворюємо модель завдання у доменну модель і повертаємо.
}

// DeleteByID видаляє завдання за його ID.
// Якщо сталася помилка під час видалення, повертається ця помилка.
func (r taskRepository) DeleteByID(id uint64) error {
	err := r.coll.Find(db.Cond{"id": id}).Delete() // Виконуємо видалення завдання за його ID.
	if err != nil {
		return err // Якщо сталася помилка, повертаємо її.
	}
	return nil // Якщо все пройшло успішно, повертаємо nil.
}

// UpdateStatus оновлює статус завдання за його ID.
// Якщо сталася помилка під час оновлення, повертається ця помилка.
func (r taskRepository) UpdateStatus(id uint64, status domain.TaskStatus) error {
	err := r.coll.
		Find(db.Cond{"id": id}). // Знаходимо завдання за його ID.
		Update(map[string]interface{}{
			"status": status, // Оновлюємо статус завдання.
		})
	if err != nil { // Якщо сталася помилка, повертаємо її.
		return err
	}
	return nil // Якщо все пройшло успішно, повертаємо nil.
}
