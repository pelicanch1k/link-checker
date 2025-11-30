package repository

import (
	"sync"

	"github.com/pelicanch1k/link-checker/internal/domain"
)

type TaskRepository interface {
	Save(task *domain.Task) error
	FindByIDs(ids []int) ([]*domain.Task, error)
	GetNextID() int
}

type InMemoryTaskRepository struct {
	mu      sync.RWMutex
	tasks   map[int]*domain.Task
	counter int
}

func NewInMemoryTaskRepository() *InMemoryTaskRepository {
	return &InMemoryTaskRepository{
		tasks:   make(map[int]*domain.Task),
		counter: 0,
	}
}

func (r *InMemoryTaskRepository) Save(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = task
	return nil
}

func (r *InMemoryTaskRepository) FindByIDs(ids []int) ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	var result []*domain.Task
	for _, id := range ids {
		if task, exists := r.tasks[id]; exists {
			result = append(result, task)
		}
	}
	
	return result, nil
}

func (r *InMemoryTaskRepository) GetNextID() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counter++
	return r.counter
}