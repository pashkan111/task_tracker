package api

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"strconv"
	"task_tracker/src/entities"
	"task_tracker/src/errors/api_errors"
	"task_tracker/src/services"
	"task_tracker/src/utils"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/sirupsen/logrus"
)

func InitUserRoutes(router *mux.Router, pool *pgxpool.Pool, log *logrus.Logger) {
	router.HandleFunc("/users", CreateUser(pool, log)).Methods("POST")
	router.HandleFunc("/users/{userId}", UpdateUser(pool, log)).Methods("PATCH")
	router.HandleFunc("/users/{userId}", DeleteUser(pool, log)).Methods("DELETE")
	router.HandleFunc("/users", GetUsers(pool, log)).Methods("GET")
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

		created_user, err := services.CreateUser(
			r.Context(),
			pool,
			log,
			*user_data_validated,
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

		updated_user, err := services.UpdateUser(
			r.Context(),
			pool,
			log,
			*validated_user_data,
			user_id,
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

		response := entities.UserUpdateResponse(updated_user)
		json.NewEncoder(w).Encode(response)
	}
}

func DeleteUser(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
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

		err = services.DeleteUser(r.Context(), pool, log, user_id)
		if err != nil {
			resp := entities.ErrorResponse{Error: err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func GetUsers(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")

		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			page = 1
		}

		users_per_page := 5

		users, users_count, err := services.GetUsers(
			r.Context(),
			pool,
			log,
			(page-1)*users_per_page,
			users_per_page,
		)
		if err != nil {
			resp := entities.ErrorResponse{Error: err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}

		response := entities.GetUsersResponse{
			Users:    users,
			Page:     page,
			LastPage: int(math.Ceil(float64(users_count) / float64(users_per_page))),
		}
		json.NewEncoder(w).Encode(response)
	}
}
