package entities

import "time"

type CreateTaskRequest struct {
	TaskName string `json:"taskName" validate:"required"`
	UserId   int    `json:"userId" validate:"required"`
}

type CreateTaskResponse struct {
	TaskId    int       `json:"taskId"`
	TaskName  string    `json:"taskName"`
	UserId    int       `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}
