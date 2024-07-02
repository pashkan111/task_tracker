package repository

import (
	"context"
	"task_tracker/src/entities"
	"task_tracker/src/errors/repo_errors"

	"github.com/jackc/pgconn"

	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

func CreateUser(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	user entities.User,
) (int, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return 0, err
	}
	defer conn.Release()

	var userID int
	err_create := conn.QueryRow(
		ctx,
		`INSERT INTO users (passport_serie, passport_number, surname, name) 
		VALUES ($1, $2, $3, $4) 
		RETURNING user_id`,
		user.PassportSerie, user.PassportNumber, user.Surname, user.Name,
	).Scan(&userID)

	if err_create != nil {
		var pg_err *pgconn.PgError
		if errors.As(err_create, &pg_err) {
			if pg_err.Code == "23505" {
				log.Errorf("error: %s. Detail: %s", pg_err.Error(), pg_err.Detail)
				return 0, repo_errors.ObjectAlreadyExistsError{}
			}
		} else {
			log.Error("Error creating user:", err_create)
			return 0, repo_errors.OperationError{}
		}
	}
	return userID, err
}

func UpdateUser(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	user entities.User,
	user_id int,
) error {
	conn, err := pool.Acquire(ctx)

	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return err
	}
	defer conn.Release()

	_, err_exec := conn.Exec(
		ctx,
		`UPDATE users 
		SET passport_serie=$1, passport_number=$2, surname=$3, name=$4 
		WHERE user_id=$5
		RETURNING user_id
		`,
		user.PassportSerie, user.PassportNumber, user.Surname, user.Name, user_id,
	)
	return err_exec
}

func DeleteUser(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	user_id int,
) error {
	conn, err := pool.Acquire(ctx)

	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(
		ctx,
		`DELETE FROM users 
		WHERE user_id=$1`,
		user_id,
	)
	if err != nil {
		log.Error("Error deleting user:", err)
	}

	return err
}

func GetUsers(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	offset int,
	limit int,
) ([]entities.User, error) {
	conn, err := pool.Acquire(ctx)

	var users []entities.User
	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return users, err
	}
	defer conn.Release()

	rows, err := conn.Query(
		ctx,
		`SELECT user_id, passport_serie, passport_number, surname, name 
		FROM users
		LIMIT $1
		OFFSET $2;
		`,
		limit, offset,
	)
	if err != nil {
		log.Error("Error selecting users:", err)
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entities.User
		err = rows.Scan(&user.Id, &user.PassportSerie, &user.PassportNumber, &user.Surname, &user.Name)
		if err != nil {
			log.Error("Error scanning user:", err)
			return users, err
		}
		users = append(users, user)
	}

	return users, err
}
