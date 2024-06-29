package main

import (
	"context"
	"fmt"
	"task_tracker/src/api"
	"task_tracker/src/repository"
	"task_tracker/src/utils"
	"time"

	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	log := utils.GetLogger()
	postgres_pool := repository.GetPostgresPool(ctx, log)

	router := mux.NewRouter()
	api.InitUserRoutes(router, postgres_pool, log)

	fmt.Println("Server is running on port 8080")
	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
