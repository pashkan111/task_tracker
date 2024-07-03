package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"task_tracker/src/entities"
	"task_tracker/src/services"
	"task_tracker/src/utils"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/sirupsen/logrus"
)

func InitUserRoutes(router *mux.Router, pool *pgxpool.Pool, log *logrus.Logger) {
	router.HandleFunc("/users", CreateUser(pool, log)).Methods("POST")
	router.HandleFunc("/users/{userId}", UpdateUser(pool, log)).Methods("PATCH")
}

func CreateUser(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var user entities.UserCreateRequest
		user_data_validated, err := utils.ValidateRequestData(user, r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp := entities.ErrorResponse{Error: err.Error()}
			json.NewEncoder(w).Encode(resp)
			return
		}

		created_user, err := services.CreateUser(r.Context(), pool, log, *user_data_validated)
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

func UpdateUser(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		userID := vars["userId"]

		if userID == "" {
			resp := entities.ErrorResponse{Error: "Parametr userId is required"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
		user_id, err := strconv.Atoi(userID)
		if err != nil {
			resp := entities.ErrorResponse{Error: "Parametr userId must be a number"}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		validated_user_data, err_parse := utils.ValidateRequestData(entities.UserUpdateRequest{}, r.Body)

		if err_parse != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp := entities.ErrorResponse{Error: err_parse.Error()}
			json.NewEncoder(w).Encode(resp)
			return
		}

		updated_user, err := services.UpdateUser(r.Context(), pool, log, *validated_user_data, user_id)
		if err != nil {
			resp := entities.ErrorResponse{Error: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		response := entities.UserUpdateResponse(updated_user)
		json.NewEncoder(w).Encode(response)
	}
}
