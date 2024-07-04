package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"task_tracker/src/entities"
	"task_tracker/src/errors/api_errors"
	"task_tracker/src/services"
	"task_tracker/src/utils"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/sirupsen/logrus"
)

func InitTaskRoutes(router *mux.Router, pool *pgxpool.Pool, log *logrus.Logger) {
	router.HandleFunc("/tasks", createTask(pool, log)).Methods("POST")
}

func createTask(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var task entities.CreateTaskRequest
		task_data_validated, err := utils.ValidateRequestData(task, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp := entities.ErrorResponse{Error: err.Error()}
			json.NewEncoder(w).Encode(resp)
			return
		}

		created_task, err := services.CreateTask(
			r.Context(),
			pool,
			log,
			*task_data_validated,
		)
		if err != nil {
			resp := entities.ErrorResponse{Error: err.Error()}
			var bad_request_error *api_errors.BadRequestError

			if errors.As(err, &bad_request_error) {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		json.NewEncoder(w).Encode(created_task)
	}
}
