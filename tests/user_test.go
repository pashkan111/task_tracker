package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"task_tracker/src/api"
	"task_tracker/src/entities"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestRegisterUserHandler(t *testing.T) {
	pool, cleanup, err := SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	log := SetupLogger()

	router := mux.NewRouter()
	api.InitUserRoutes(router, pool, log)

	passport_serie := 2233
	passport_number := 895044
	name := "Ivan"
	surname := "Petrov"

	user := entities.UserCreateRequest{
		PassportNumber: fmt.Sprintf("%d %d", passport_serie, passport_number),
		Name:           name,
		Surname:        surname,
	}
	body, _ := json.Marshal(user)

	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response entities.UserCreateResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)
	assert.Equal(t, response.Name, name)
	assert.Equal(t, response.Surname, surname)
	assert.Equal(t, response.PassportNumber, passport_number)
	assert.Equal(t, response.PassportSerie, passport_serie)

	var userFromDB entities.User
	row := pool.QueryRow(
		context.Background(),
		"SELECT * FROM users WHERE surname = $1",
		surname,
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
