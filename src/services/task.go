package services

import (
	"context"
	"errors"
	"fmt"

	"task_tracker/src/entities"
	"task_tracker/src/errors/api_errors"
	"task_tracker/src/errors/repo_errors"
	"task_tracker/src/repository"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/sirupsen/logrus"
)

func CreateTask(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	task entities.CreateTaskRequest,
) (*entities.CreateTaskResponse, error) {

	created_task, err := repository.CreateTask(
		ctx, pool, log, task,
	)
	if err != nil {
		if errors.Is(err, repo_errors.ObjectNotFoundError{}) {
			return nil, &api_errors.BadRequestError{
				Detail: fmt.Sprintf("User with id=%d does not exist", task.UserId),
			}
		}
		return nil, &api_errors.InternalServerError{}
	}
	return created_task, nil
}

func FinishTask(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	task_id int,
) error {
	err := repository.FinishTask(
		ctx, pool, log, task_id,
	)
	if err != nil {
		return &api_errors.InternalServerError{}
	}
	return nil
}
