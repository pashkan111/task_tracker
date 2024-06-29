package services

import (
	"context"
	"fmt"
	"strings"

	"strconv"
	"task_tracker/src/entities"
	"task_tracker/src/repository"

	"github.com/jackc/pgx/v5/pgxpool"
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
		return entities.User{}, fmt.Errorf("Incorrect passportNumber format. Pass passport serie and passport number devided with space")
	}
	passport_serie, err_serie := strconv.Atoi(passport_data[0])
	passport_number, err_num := strconv.Atoi(passport_data[1])
	if err_serie != nil || err_num != nil {
		return entities.User{}, fmt.Errorf("Incorrect passportNumber format. Passport serie and passport number must be numbers, not strings")
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
		return entities.User{}, fmt.Errorf("Error when creating user. Err: %s", err.Error())
	}
	user_to_create.Id = user_id
	return user_to_create, nil
}
