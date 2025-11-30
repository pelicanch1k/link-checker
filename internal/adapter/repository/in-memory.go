package repository

import (
	"github.com/pelicanch1k/link-checker/internal/domain"
	
	"sync"
	"fmt"
	"time"
)

type TaskRepository interface {
	Create(task *domain.Task) error
	FindByID(id string) (*domain.Task, error)
	Update(task *domain.Task) error
	FindPending() ([]*domain.Task, error)
	GetAll() []*domain.Task
}

type InMemoryTaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks: make(map[string]*domain.Task),
	}
}

func (r *InMemoryTaskRepository) Create(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = task
	return nil
}

func (r *InMemoryTaskRepository) FindByID(id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	task, ok := r.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found")
	}
	return task, nil
}

func (r *InMemoryTaskRepository) Update(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	task.UpdatedAt = time.Now()
	r.tasks[task.ID] = task
	return nil
}

func (r *InMemoryTaskRepository) FindPending() ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var pending []*domain.Task
	for _, task := range r.tasks {
		if task.Status == domain.TaskPending || task.Status == domain.TaskProcessing {
			pending = append(pending, task)
		}
	}
	return pending, nil
}

func (r *InMemoryTaskRepository) GetAll() []*domain.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var tasks []*domain.Task
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}