package repository

import (
	"context"
	"fmt"
	"strings"
	"task_tracker/src/entities"
	"task_tracker/src/errors/repo_errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"

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
			log.Error("Error creating user: ", err_create)
			return 0, repo_errors.OperationError{}
		}
	}
	return userID, err
}

func UpdateUser(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	user entities.UserUpdateRepo,
	user_id int,
) (entities.User, error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Error("Error with acquiring connection:", err)
		return entities.User{}, err
	}
	defer conn.Release()

	set_clauses := []string{}
	args := []interface{}{}
	argID := 1

	if user.PassportSerie != nil {
		set_clauses = append(set_clauses, fmt.Sprintf("passport_serie=$%d", argID))
		args = append(args, *user.PassportSerie)
		argID++
	}
	if user.PassportNumber != nil {
		set_clauses = append(set_clauses, fmt.Sprintf("passport_number=$%d", argID))
		args = append(args, *user.PassportNumber)
		argID++
	}
	if user.Surname != nil {
		set_clauses = append(set_clauses, fmt.Sprintf("surname=$%d", argID))
		args = append(args, *user.Surname)
		argID++
	}
	if user.Name != nil {
		set_clauses = append(set_clauses, fmt.Sprintf("name=$%d", argID))
		args = append(args, *user.Name)
		argID++
	}

	if len(set_clauses) == 0 {
		return entities.User{}, fmt.Errorf("No fields to update")
	}

	set_query := strings.Join(set_clauses, ", ")
	query := fmt.Sprintf(
		`UPDATE users 
		SET %s 
		WHERE user_id=$%d 
		RETURNING user_id, passport_serie, passport_number, surname, name;`,
		set_query,
		argID,
	)
	args = append(args, user_id)

	var updated_user entities.User
	err = conn.QueryRow(ctx, query, args...).Scan(
		&updated_user.Id,
		&updated_user.PassportSerie,
		&updated_user.PassportNumber,
		&updated_user.Surname,
		&updated_user.Name,
	)

	if err != nil {
		var pg_err *pgconn.PgError

		if errors.As(err, &pg_err) {
			if pg_err.Code == "23505" {
				log.Errorf("error: %s. Detail: %s", pg_err.Error(), pg_err.Detail)
				return entities.User{}, repo_errors.ObjectAlreadyExistsError{}
			}
		} else if err.Error() == pgx.ErrNoRows.Error() {
			log.Errorf("error: %s. Detail: %s=%d", err.Error(), "user_id", user_id)
			return entities.User{}, repo_errors.ObjectNotFoundError{}
		} else {
			log.Errorf("Error updating user: %s", err)
			return entities.User{}, repo_errors.OperationError{}
		}
	}

	if updated_user.Id == 0 {
		return entities.User{}, repo_errors.ObjectNotFoundError{}
	}
	return updated_user, nil
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
		return repo_errors.OperationError{}
	}

	return nil
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
