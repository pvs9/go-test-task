package service

import (
	"github.com/todo"
	"github.com/todo/pkg/queue"
	"github.com/todo/pkg/repository"
)

type TodoItem interface {
	Create(item todo.TodoItem) (int64, error)
	GetLatest() (*todo.TodoItem, error)
}

type Service struct {
	TodoItem
}

func NewService(repositories *repository.Repository, queues *queue.Queue) *Service {
	return &Service{
		TodoItem: NewTodoItemService(queues.Publisher, repositories.TodoItem),
	}
}
