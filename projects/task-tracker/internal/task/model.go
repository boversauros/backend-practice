package task

import "time"

type Status string

const (
    StatusTodo Status = "todo"
    StatusInProgress Status = "in-progress"
    StatusDone Status = "done"
)

type Task struct {
    ID          int       `json:"id"`
    Description string    `json:"description"`
    Status      Status    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}