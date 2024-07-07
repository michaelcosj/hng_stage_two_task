package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	database "github.com/michaelcosj/hng-task-two/internal/db"
	"github.com/michaelcosj/hng-task-two/internal/service"
)

func New(db database.Db) *http.Server {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 6969
	}

	querier := database.New(db)
	repo := database.NewRepoQuerier(querier, db)

	svc := service.New(repo)
	handler := RegisterRoutes(svc)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      handler,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return httpServer
}
