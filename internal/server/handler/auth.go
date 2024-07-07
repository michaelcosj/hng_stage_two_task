package handler

import (
	"encoding/json"
	"net/http"

	"github.com/michaelcosj/hng-task-two/internal/app"
	"github.com/michaelcosj/hng-task-two/internal/service"
)

func (s *Handler) AuthRegister(w http.ResponseWriter, r *http.Request) error {
	var req RegisterUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return app.InvalidJson()
	}

	if problems := req.Validate(); len(problems) > 0 {
		return app.NewValidationError(problems)
	}

	data, err := s.service.Register(r.Context(), service.RegisterParams{
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Password:  req.Password,
		Phone:     req.Phone,
	})

	if err != nil {
		return app.ApiErrorFrom(err)
	}

	writeJSON(w, http.StatusCreated, SuccessResponse{
		Status:  "success",
		Message: "Registration successful",
		Data:    data,
	})

	return nil
}

func (s *Handler) AuthLogin(w http.ResponseWriter, r *http.Request) error {
	var req LoginUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return app.InvalidJson()
	}

	if problems := req.Validate(); len(problems) > 0 {
		return app.NewValidationError(problems)
	}

	data, err := s.service.Login(r.Context(), service.LoginParams{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		return app.ApiErrorFrom(err)
	}

	writeJSON(w, http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "Login successful",
		Data:    data,
	})

	return nil
}
