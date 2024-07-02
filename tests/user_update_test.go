package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task_tracker/src/api"
	"task_tracker/src/entities"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestUpdateUserHandler__OK(t *testing.T) {
	pool, cleanup, err := SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	log := SetupLogger()

	router := mux.NewRouter()
	api.InitUserRoutes(router, pool, log)

	_, err_create := pool.Exec(
		context.Background(),
		`INSERT INTO users (user_id, passport_serie, passport_number, surname, name)
		VALUES (1, 1212, 232323, 'Ivanov', 'Ivan');
		`,
	)
	require.NoError(t, err_create)

	passport_serie := 2233
	passport_number := 895044
	name := "Ivan"
	surname := "Petrov"

	user := entities.UserUpdateRequest{
		PassportSerie:  passport_serie,
		PassportNumber: passport_number,
		Surname:        surname,
		Name:           name,
	}
	body, _ := json.Marshal(user)

	req, err := http.NewRequest("PATCH", "/users/1", bytes.NewBuffer(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response entities.UserUpdateResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, response.Name, name)
	assert.Equal(t, response.Surname, surname)
	assert.Equal(t, response.PassportNumber, passport_number)
	assert.Equal(t, response.PassportSerie, passport_serie)

	var userFromDB entities.User
	row := pool.QueryRow(
		context.Background(),
		"SELECT * FROM users WHERE user_id = 1;",
	)

	err = row.Scan(
		&userFromDB.Id,
		&userFromDB.PassportSerie,
		&userFromDB.PassportNumber,
		&userFromDB.Surname,
		&userFromDB.Name,
	)
	require.NoError(t, err)

	assert.Equal(t, name, userFromDB.Name)
	assert.Equal(t, passport_number, userFromDB.PassportNumber)
}

func TestUpdateUserHandler__BadRequest(t *testing.T) {
	pool, cleanup, err := SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	log := SetupLogger()

	router := mux.NewRouter()
	api.InitUserRoutes(router, pool, log)

	_, err_create := pool.Exec(
		context.Background(),
		`INSERT INTO users (user_id, passport_serie, passport_number, surname, name)
		VALUES (1, 1212, 232323, 'Ivanov', 'Ivan');
		`,
	)
	pool.Exec(
		context.Background(),
		`INSERT INTO users (user_id, passport_serie, passport_number, surname, name)
		VALUES (1, 1212, 895044, 'Ivanov', 'Ivan');
		`,
	)
	require.NoError(t, err_create)

	passport_serie := 2233
	passport_number := 895044
	name := "Ivan"
	surname := "Petrov"

	user := entities.UserUpdateRequest{
		PassportSerie:  passport_serie,
		PassportNumber: passport_number,
		Surname:        surname,
		Name:           name,
	}
	body, _ := json.Marshal(user)

	req, err := http.NewRequest("PATCH", "/users/12", bytes.NewBuffer(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response entities.ErrorResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, "User not found. User_id: 12", response.Error)

	var userFromDB entities.User
	row := pool.QueryRow(
		context.Background(),
		"SELECT * FROM users WHERE user_id = 1;",
	)

	err = row.Scan(
		&userFromDB.Id,
		&userFromDB.PassportSerie,
		&userFromDB.PassportNumber,
		&userFromDB.Surname,
		&userFromDB.Name,
	)
	require.NoError(t, err)

	assert.Equal(t, name, userFromDB.Name)
	assert.Equal(t, passport_number, userFromDB.PassportNumber)
}
