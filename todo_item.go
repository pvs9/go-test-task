package todo

import "time"

type TodoItem struct {
	Id          int64     `json:"id" db:"id"`
	Description string    `json:"description" db:"description"`
	DueDate     time.Time `json:"due_date" db:"due_date"`
}
