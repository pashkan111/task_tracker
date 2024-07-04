package tests

import (
	"context"
	"encoding/json"
	"log"
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

func createUsers(pool *pgxpool.Pool) {
	queries := []string{
		`INSERT INTO users (user_id, passport_serie, passport_number, surname, name) 
			VALUES (1, 1212, 232323, 'Ivanov', 'Ivan');`,
		`INSERT INTO users (user_id, passport_serie, passport_number, surname, name) 
			VALUES (2, 3434, 454545, 'Petrov', 'Petr');`,
		`INSERT INTO users (user_id, passport_serie, passport_number, surname, name) 
			VALUES (3, 5656, 676767, 'Sidorov', 'Sidr');`,
		`INSERT INTO users (user_id, passport_serie, passport_number, surname, name) 
			VALUES (4, 7878, 898989, 'Smirnov', 'Sergey');`,
		`INSERT INTO users (user_id, passport_serie, passport_number, surname, name) 
			VALUES (5, 9090, 111111, 'Kuznetsov', 'Alexey');`,
		`INSERT INTO users (user_id, passport_serie, passport_number, surname, name) 
			VALUES (6, 1213, 121212, 'Popov', 'Dmitry');`,
		`INSERT INTO users (user_id, passport_serie, passport_number, surname, name) 
			VALUES (7, 1414, 232424, 'Vasilev', 'Vladimir');`,
		`INSERT INTO users (user_id, passport_serie, passport_number, surname, name) 
			VALUES (8, 1616, 343535, 'Mikhailov', 'Mikhail');`,
		`INSERT INTO users (user_id, passport_serie, passport_number, surname, name) 
			VALUES (9, 1818, 454646, 'Fedorov', 'Fedor');`,
	}

	for _, query := range queries {
		_, err := pool.Exec(context.Background(), query)
		if err != nil {
			log.Fatalf("Failed to execute query: %v\n", err)
		}
	}
}

func TestGetUsersHandler__OK(t *testing.T) {
	pool, cleanup, err := SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	log := SetupLogger()

	router := mux.NewRouter()
	api.InitUserRoutes(router, pool, log)

	createUsers(pool)

	req, err := http.NewRequest("GET", "/users", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response entities.GetUsersResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	require.NoError(t, err)

	assert.Equal(t, 5, len(response.Users))
	assert.Equal(t, 1, response.Page)
	assert.Equal(t, 2, response.LastPage)
}
