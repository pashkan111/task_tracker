package tests

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"task_tracker/src/api"
	"task_tracker/src/entities"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/require"
)

func createUsersAndTasks(pool *pgxpool.Pool) {
	pool.Exec(context.Background(),
		`INSERT INTO users (user_id, passport_serie, passport_number, surname, name) 
			VALUES (1, 1212, 232323, 'Ivanov', 'Ivan');`,
	)
	pool.Exec(
		context.Background(),
		`INSERT INTO tasks (task_id, user_id, name, created_at, due_date)
		VALUES (1, 1, 'task1', '2024-07-09 11:50:07', '2024-07-09 19:24:07');
		`,
	)
	pool.Exec(
		context.Background(),
		`INSERT INTO tasks (task_id, user_id, name, created_at, due_date)
		VALUES (2, 1, 'task2', '2024-07-06 10:55:07', '2024-07-08 14:24:07');
		`,
	)
	pool.Exec(
		context.Background(),
		`INSERT INTO tasks (task_id, user_id, name, created_at, due_date)
		VALUES (3, 1, 'task3', '2024-07-10 11:50:07', '2024-07-11 15:24:07');
		`,
	)
}

func TestGetUserActivitiesHandler__OK(t *testing.T) {
	pool, cleanup, err := SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	log := SetupLogger()

	router := mux.NewRouter()
	api.InitUserRoutes(router, pool, log)

	createUsersAndTasks(pool)

	req, err := http.NewRequest("GET", "/user-activities/1", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response entities.UserActivityResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, 3, len(response.Tasks))
	assert.Equal(t, "task1", response.Tasks[1].TaskName)
	assert.Equal(t, "task3", response.Tasks[0].TaskName)
}
