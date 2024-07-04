package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"task_tracker/src/entities"
	"task_tracker/src/errors/api_errors"
	"task_tracker/src/services"
	"task_tracker/src/utils"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/sirupsen/logrus"
)

func InitUserRoutes(router *mux.Router, pool *pgxpool.Pool, log *logrus.Logger) {
	router.HandleFunc("/users", createUser(pool, log)).Methods("POST")
	router.HandleFunc("/users/{userId}", updateUser(pool, log)).Methods("PATCH")
	router.HandleFunc("/users/{userId}", deleteUser(pool, log)).Methods("DELETE")
	router.HandleFunc("/users", getUsers(pool, log)).Methods("GET")
	router.HandleFunc("/user-activities/{userId}", getUserActivities(pool, log)).Methods("GET")
}

func createUser(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
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

func updateUser(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		userID := vars["userId"]

		if userID == "" {
			resp := entities.ErrorResponse{Error: api_errors.BadRequestError{Detail: "Parametr userId is required"}.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
		user_id, err := strconv.Atoi(userID)
		if err != nil {
			resp := entities.ErrorResponse{Error: api_errors.BadRequestError{Detail: "Parametr userId must be a number"}.Error()}
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

func deleteUser(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		userID := vars["userId"]

		if userID == "" {
			resp := entities.ErrorResponse{Error: api_errors.BadRequestError{Detail: "Parametr userId is required"}.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
		user_id, err := strconv.Atoi(userID)
		if err != nil {
			resp := entities.ErrorResponse{Error: api_errors.BadRequestError{Detail: "Parametr userId must be a number"}.Error()}
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

func getUsers(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
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

func getUserActivities(pool *pgxpool.Pool, log *logrus.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		userID := vars["userId"]

		if userID == "" {
			resp := entities.ErrorResponse{Error: api_errors.BadRequestError{Detail: "Parametr userId is required"}.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}
		user_id, err := strconv.Atoi(userID)
		if err != nil {
			resp := entities.ErrorResponse{Error: api_errors.BadRequestError{Detail: "Parametr userId must be a number"}.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			return
		}

		date_from := r.URL.Query().Get("dateFrom")
		date_to := r.URL.Query().Get("dateTo")

		date_time_format := "YYYY-MM-DD HH:MM"

		var date_from_parsed *time.Time
		var date_to_parsed *time.Time

		if date_from != "" {
			*date_from_parsed, err = time.Parse(date_from, date_time_format)
			if err != nil {
				resp := entities.ErrorResponse{
					Error: api_errors.BadRequestError{
						Detail: fmt.Sprintf("Parametr dateFrom must be a date. Format %s", date_time_format),
					}.Error(),
				}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(resp)
				return
			}
		}

		if date_to != "" {
			*date_to_parsed, err = time.Parse(date_to, date_time_format)
			if err != nil {
				resp := entities.ErrorResponse{
					Error: api_errors.BadRequestError{
						Detail: fmt.Sprintf("Parametr dateTo must be a date. Format %s", date_time_format),
					}.Error(),
				}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(resp)
				return
			}
		}

		user_activity_filters := entities.UserActivityRequest{
			UserId:   user_id,
			DateFrom: date_from_parsed,
			DateTo:   date_to_parsed,
		}
		user_activities, err := services.GetUserActivities(
			r.Context(),
			pool,
			log,
			user_activity_filters,
		)
		if err != nil {
			resp := entities.ErrorResponse{Error: err.Error()}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
		resp := entities.UserActivityResponse{
			UserId: user_id,
			Tasks:  user_activities,
		}
		json.NewEncoder(w).Encode(resp)
	}
}
