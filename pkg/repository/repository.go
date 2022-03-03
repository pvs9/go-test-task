package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/todo"
)

type TodoItem interface {
	Create(item todo.TodoItem) (int64, error)
	GetLatest() (*todo.TodoItem, error)
}

type Repository struct {
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		TodoItem: NewTodoItemMySQL(db),
	}
}
