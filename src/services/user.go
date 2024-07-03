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
	passport_data, err := getPassportDataFromString(user.PassportNumber)
	if err != nil {
		return entities.User{}, api_errors.BadRequestError{
			Detail: err.Error(),
		}
	}
	passport_serie := passport_data[0]
	passport_number := passport_data[1]

	user_to_create := entities.User{
		PassportSerie:  *passport_serie,
		PassportNumber: *passport_number,
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
	passport_data := []*int{nil, nil}
	var err error

	if user.PassportNumber != nil {
		passport_data, err = getPassportDataFromString(*user.PassportNumber)
		if err != nil {
			return entities.User{}, api_errors.BadRequestError{
				Detail: err.Error(),
			}
		}
	}
	passport_serie := passport_data[0]
	passport_number := passport_data[1]

	user_to_update := entities.UserUpdateRepo{
		PassportSerie:  passport_serie,
		PassportNumber: passport_number,
		Surname:        user.Surname,
		Name:           user.Name,
	}
	updated_user, err := repository.UpdateUser(
		ctx, pool, log, user_to_update, user_id,
	)
	if err != nil {
		if errors.Is(err, repo_errors.ObjectAlreadyExistsError{}) {
			return entities.User{}, api_errors.BadRequestError{
				Detail: fmt.Sprintf("%s. userId=%d", err, user_id),
			}
		} else if errors.Is(err, repo_errors.ObjectNotFoundError{}) {
			return entities.User{}, api_errors.BadRequestError{
				Detail: fmt.Sprintf("%s. userId=%d", err, user_id),
			}
		}
		return entities.User{}, api_errors.BadRequestError{Detail: err.Error()}
	}
	return updated_user, nil
}

func getPassportDataFromString(passport_number_string string) ([]*int, error) {
	passport_data := strings.Split(passport_number_string, " ")
	if len(passport_data) < 2 {
		return []*int{}, fmt.Errorf("Incorrect passportNumber format. Passport serie and passport number should be devided with space")
	}
	passport_serie, err_serie := strconv.Atoi(passport_data[0])
	passport_number, err_num := strconv.Atoi(passport_data[1])
	if err_serie != nil || err_num != nil {
		return []*int{}, fmt.Errorf("Incorrect passportNumber format. Passport serie and passport number must be numbers, not strings")
	}
	return []*int{&passport_serie, &passport_number}, nil
}
