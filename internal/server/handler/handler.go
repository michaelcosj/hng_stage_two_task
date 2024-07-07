package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/michaelcosj/hng-task-two/internal/app"
	"github.com/michaelcosj/hng-task-two/internal/service"
)

type ApiFunc func(w http.ResponseWriter, r *http.Request) error

type Handler struct {
	service service.Service
}

func New(svc service.Service) *Handler {
	return &Handler{service: svc}
}

func Handle(handler ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			if apiError, ok := err.(app.ApiError); ok {
				writeJSON(w, apiError.StatusCode, apiError)
			} else if apiValidationError, ok := err.(app.ApiValidationError); ok {
				writeJSON(w, http.StatusUnprocessableEntity, apiValidationError)
			} else {
				errResp := map[string]any{
					"status":     "internal server error",
					"statusCode": http.StatusInternalServerError,
					"message":    "something went wrong, please don't fail me",
				}

				writeJSON(w, http.StatusInternalServerError, errResp)
			}

			log.Printf("an error occured: %v", err)
		}
	}
}

type SuccessResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func getAuthUserFromContext(ctx context.Context) (uuid.UUID, error) {
	userIdValue, ok := ctx.Value("userId").(string)
	if !ok {
		return uuid.UUID{}, app.ErrAuthenticationFailed
	}

	userId, err := uuid.Parse(userIdValue)
	if err != nil {
		return uuid.UUID{}, app.ApiErrorFrom(errors.Join(err, app.ErrAuthenticationFailed))
	}

	return userId, nil
}
