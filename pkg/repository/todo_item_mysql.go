package repository

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/todo"
)

type TodoItemMySQL struct {
	db *sqlx.DB
}

func NewTodoItemMySQL(db *sqlx.DB) *TodoItemMySQL {
	return &TodoItemMySQL{db: db}
}

func (r *TodoItemMySQL) Create(item todo.TodoItem) (int64, error) {
	createItemQuery := fmt.Sprintf("INSERT INTO %s (description, due_date) values (?, ?)", todoItemsTable)

	result, err := r.db.Exec(createItemQuery, item.Description, item.DueDate)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *TodoItemMySQL) GetLatest() (*todo.TodoItem, error) {
	var item todo.TodoItem
	query := fmt.Sprintf(`SELECT ti.id, ti.description, ti.due_date FROM %s ti ORDER BY due_date DESC LIMIT 1`,
		todoItemsTable)

	if err := r.db.Get(&item, query); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return &item, err
	}

	return &item, nil
}
