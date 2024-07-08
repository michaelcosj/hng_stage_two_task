package handler

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/michaelcosj/hng-task-two/internal/app"
)

func (s *Handler) GetUser(w http.ResponseWriter, r *http.Request) error {
	authUserId, err := getAuthUserFromContext(r.Context())
	if err != nil {
		return err
	}

	userIdStr := r.PathValue("userId")
	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return app.InvalidRequestData(fmt.Errorf("error parsing uuid: %v", err))
	}

	data, err := s.service.GetUser(r.Context(), authUserId, userId)
	if err != nil {
		return app.ApiErrorFrom(err)
	}

	writeJSON(w, http.StatusCreated, SuccessResponse{
		Status:  "success",
		Message: "User found successfully",
		Data:    data,
	})

	return nil
}
