package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"task_tracker/src/entities"
	"task_tracker/src/services"
	"task_tracker/src/utils"

	"github.com/go-playground/validator/v10"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/sirupsen/logrus"
)

func InitUserRoutes(router *mux.Router, pool *pgxpool.Pool, log *logrus.Logger) {
	router.HandleFunc("/users", CreateUser(pool, log)).Methods("POST")
	router.HandleFunc("/users/{userId}", UpdateUser(pool, log)).Methods("PUT")
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
			w.WriteHeader(http.StatusBadRequest)

			for _, err := range err.(validator.ValidationErrors) {
				fieldError := fmt.Sprintf("Validation failed on field '%s', condition: '%s'", err.Field(), err.Tag())
				json.NewEncoder(w).Encode(entities.ErrorResponse{Error: fieldError})
				return
			}

			resp := entities.ErrorResponse{Error: fmt.Sprintf("Validation error. %s", err.Error())}
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

		var user entities.UserUpdateRequest
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&user); err != nil {
			resp := entities.ErrorResponse{Error: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		updated_user, err := services.UpdateUser(r.Context(), pool, log, user, user_id)
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
