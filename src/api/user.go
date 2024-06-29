package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"task_tracker/src/entities"
	"task_tracker/src/services"
	"task_tracker/src/utils"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

func InitUserRoutes(router *mux.Router, pool *pgxpool.Pool, log *logrus.Logger) {
	router.HandleFunc("/users", CreateUser(pool, log)).Methods("POST")
}

func CreateUser(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var user entities.UserCreateRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			resp := entities.ErrorResponse{Error: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		if err := utils.ValidateData(&user); err != nil {
			resp := entities.ErrorResponse{Error: fmt.Sprintf("Validation error. %s", err.Error())}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		created_user, err := services.CreateUser(r.Context(), pool, log, user)
		if err != nil {
			resp := entities.ErrorResponse{Error: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		response := entities.UserCreateResponse(created_user)
		json.NewEncoder(w).Encode(response)
	}
}
