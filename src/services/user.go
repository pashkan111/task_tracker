package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"task_tracker/src/entities"
	"task_tracker/src/errors/api_errors"
	"task_tracker/src/errors/repo_errors"
	"task_tracker/src/repository"

	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/sirupsen/logrus"
)

func CreateUser(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	user entities.UserCreateRequest,
) (entities.User, error) {
	passport_data := strings.Split(user.PassportNumber, " ")
	if len(passport_data) < 2 {
		return entities.User{}, api_errors.BadRequestError{
			Detail: "Incorrect passportNumber format. Passport serie and passport number should be devided with space",
		}
	}
	passport_serie, err_serie := strconv.Atoi(passport_data[0])
	passport_number, err_num := strconv.Atoi(passport_data[1])
	if err_serie != nil || err_num != nil {
		return entities.User{}, api_errors.BadRequestError{
			Detail: "Incorrect passportNumber format. Passport serie and passport number must be numbers, not strings",
		}
	}

	user_to_create := entities.User{
		PassportSerie:  passport_serie,
		PassportNumber: passport_number,
		Surname:        user.Surname,
		Name:           user.Name,
	}
	user_id, err := repository.CreateUser(
		ctx, pool, log, user_to_create,
	)
	if err != nil {
		if errors.Is(err, repo_errors.ObjectAlreadyExistsError{}) {
			return entities.User{}, api_errors.BadRequestError{
				Detail: fmt.Sprintf("User with passport number %s already exists", user.PassportNumber),
			}
		}
		return entities.User{}, api_errors.InternalServerError{}
	}
	user_to_create.Id = user_id
	return user_to_create, nil
}

func UpdateUser(
	ctx context.Context,
	pool *pgxpool.Pool,
	log *logrus.Logger,
	user entities.UserUpdateRequest,
	user_id int,
) (entities.User, error) {
	user_to_update := entities.User{
		PassportSerie:  user.PassportSerie,
		PassportNumber: user.PassportNumber,
		Surname:        user.Surname,
		Name:           user.Name,
	}
	err := repository.UpdateUser(
		ctx, pool, log, user_to_update, user_id,
	)
	if err != nil {
		// if errors.As(err, &repo_errors.ObjectNotFoundError) {
		return entities.User{}, api_errors.BadRequestError{}
		// }
	}
	// return entities.User{}, fmt.Errorf("Error when updating user. Err: %s", err.Error())
	user_to_update.Id = user_id
	return user_to_update, nil
}
