package server

import (
	"net/http"

	"github.com/michaelcosj/hng-task-two/internal/server/handler"
	"github.com/michaelcosj/hng-task-two/internal/service"
)

func RegisterRoutes(svc service.Service) http.Handler {
	h := handler.New(svc)

	mux := http.NewServeMux()
	// ------ Auth Routes ------ //
	authRoutes := http.NewServeMux()
	authRoutes.HandleFunc("POST /register", handler.Handle(h.AuthRegister))
	authRoutes.HandleFunc("POST /login", handler.Handle(h.AuthLogin))

	// ------ API Routes ------ //
	apiRoutes := http.NewServeMux()
	apiRoutes.HandleFunc("GET /users/{userId}", handler.Handle(h.GetUser))
	apiRoutes.HandleFunc("GET /organisations", handler.Handle(h.GetUserOrganisations))
	apiRoutes.HandleFunc("POST /organisations", handler.Handle(h.CreateNewOrganisation))
	apiRoutes.HandleFunc("GET /organisations/{orgId}", handler.Handle(h.GetSingleOrganisation))
	apiRoutes.HandleFunc("POST /organisations/{orgId}/users", handler.Handle(h.AddUserToOrganisation))

	mux.Handle("/auth/", http.StripPrefix("/auth", authRoutes))
	mux.Handle("/api/", http.StripPrefix("/api", handler.Authenticate(apiRoutes)))

	// just incase
	mux.HandleFunc("POST /api/auth/register", handler.Handle(h.AuthRegister))
	mux.HandleFunc("POST /api/auth/login", handler.Handle(h.AuthLogin))

	return handler.Logger(handler.StripSlashes(mux))
}
