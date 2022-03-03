package service

import (
	"encoding/json"
	"github.com/todo"
	"github.com/todo/pkg/queue"
	"github.com/todo/pkg/repository"
)

type TodoItemService struct {
	publisher  queue.Publisher
	repository repository.TodoItem
}

func NewTodoItemService(publisher queue.Publisher, repository repository.TodoItem) *TodoItemService {
	return &TodoItemService{publisher: publisher, repository: repository}
}

func (s *TodoItemService) Create(item todo.TodoItem) (int64, error) {
	id, err := s.repository.Create(item)

	if err != nil {
		return 0, err
	}

	item.Id = id
	messageBody, err := json.Marshal(item)

	if err != nil {
		return 0, err
	}

	_, err = s.publisher.Publish(string(messageBody))

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *TodoItemService) GetLatest() (*todo.TodoItem, error) {
	return s.repository.GetLatest()
}
