package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"task_tracker/src/api"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func TestDeleteUserHandler__OK(t *testing.T) {
	pool, cleanup, err := SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	log := SetupLogger()

	router := mux.NewRouter()
	api.InitUserRoutes(router, pool, log)

	pool.Exec(
		context.Background(),
		`INSERT INTO users (user_id, passport_serie, passport_number, surname, name)
		VALUES (1, 1212, 232323, 'Ivanov', 'Ivan');
		`,
	)
	req, err := http.NewRequest("DELETE", "/users/1", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)

	var count_users int
	err = pool.QueryRow(
		context.Background(),
		"SELECT count(*) FROM users;",
	).Scan(&count_users)

	require.NoError(t, err)

	assert.Equal(t, 0, count_users)
}
