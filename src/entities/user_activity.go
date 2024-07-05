package entities

import "time"

type UserActivityTask struct {
	TaskID     int    `json:"task_id"`
	TaskName   string `json:"task_name"`
	Hours      int    `json:"hours"`
	Minutes    int    `json:"minutes"`
	IsFinished bool   `json:"is_finished"`
}

type UserActivityRequest struct {
	UserId   int        `json:"userId" validate:"required"`
	DateFrom *time.Time `json:"dateFrom"`
	DateTo   *time.Time `json:"dateTo"`
}

type UserActivityResponse struct {
	UserId int                `json:"userId"`
	Tasks  []UserActivityTask `json:"tasks"`
}
