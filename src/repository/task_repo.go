package repository

import (
	"context"
	"errors"
	"task_tracker/src/entities"
	"task_tracker/src/errors/repo_errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func CreateTask(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	task entities.CreateTaskRequest,
) (*entities.CreateTaskResponse, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return nil, repo_errors.OperationError{}
	}
	defer conn.Release()

	var created_task entities.CreateTaskResponse

	err_create := conn.QueryRow(
		ctx,
		`INSERT INTO tasks (user_id, task_name) 
		VALUES ($1, $2) 
		RETURNING user_id, task_name, task_id, start_time`,
		task.UserId, task.TaskName,
	).Scan(
		&created_task.UserId,
		&created_task.TaskName,
		&created_task.TaskId,
		&created_task.CreatedAt,
	)

	if err_create != nil {
		var pg_err *pgconn.PgError
		if errors.As(err_create, &pg_err) {
			if pg_err.Code == "23503" {
				log.Errorf("error: %s. Detail: %s", pg_err.Error(), pg_err.Detail)
				return nil, repo_errors.ObjectNotFoundError{}
			}
		} else {
			log.Error("Error creating task: ", err_create)
			return nil, repo_errors.OperationError{}
		}
	}
	return &created_task, nil
}

func FinishTask(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	task_id int,
) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return repo_errors.OperationError{}
	}
	defer conn.Release()

	_, err = conn.Exec(
		ctx,
		`UPDATE tasks 
		SET end_time=current_timestamp
		WHERE task_id=$1`,
		task_id,
	)

	if err != nil {
		log.Error("Error finishing task: ", err)
		return repo_errors.OperationError{}
	}
	return nil
}
