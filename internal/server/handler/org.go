package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/michaelcosj/hng-task-two/internal/app"
	"github.com/michaelcosj/hng-task-two/internal/service"
)

func (s *Handler) GetUserOrganisations(w http.ResponseWriter, r *http.Request) error {
	userId, err := getAuthUserFromContext(r.Context())
	if err != nil {
		return err
	}

	data, err := s.service.GetUserOrganisations(r.Context(), userId)
	if err != nil {
		return app.ApiErrorFrom(err)
	}

	writeJSON(w, http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "Successfully retrieved user organisations",
		Data:    data,
	})

	return nil
}

func (s *Handler) GetSingleOrganisation(w http.ResponseWriter, r *http.Request) error {
	userId, err := getAuthUserFromContext(r.Context())
	if err != nil {
		return err
	}

	orgIdStr := r.PathValue("orgId")
	orgId, err := uuid.Parse(orgIdStr)
	if err != nil {
		return app.InvalidRequestData(fmt.Errorf("error parsing uuid: %w", err))
	}

	data, err := s.service.GetUserOrganisationById(r.Context(), userId, orgId)
	if err != nil {
		return app.ApiErrorFrom(err)
	}

	writeJSON(w, http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "Successfully retrieved organisation",
		Data:    data,
	})

	return nil
}

func (s *Handler) CreateNewOrganisation(w http.ResponseWriter, r *http.Request) error {
	userId, err := getAuthUserFromContext(r.Context())
	if err != nil {
		return err
	}

	var req CreateOrgRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return app.InvalidJson()
	}

	data, err := s.service.CreateOrganisation(r.Context(), userId, service.CreateOrgParam{
		Name:        req.Name,
		Description: req.Description,
	})

	if err != nil {
		return app.ApiErrorFrom(err)
	}

	writeJSON(w, http.StatusCreated, SuccessResponse{
		Status:  "success",
		Message: "Organisation created successfully",
		Data:    data,
	})

	return nil
}

func (s *Handler) AddUserToOrganisation(w http.ResponseWriter, r *http.Request) error {
	pathVal := r.PathValue("orgId")
	orgId, err := uuid.Parse(pathVal)
	if err != nil {
		return app.InvalidRequestData(fmt.Errorf("error parsing uuid: %w", err))
	}

	var req struct {
		UserId string `json:"userId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return app.InvalidJson()
	}

	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return app.InvalidRequestData(fmt.Errorf("error parsing uuid: %v", err))
	}

	err = s.service.AddUserToOrganisation(r.Context(), orgId, userId)
	if err != nil {
		return app.ApiErrorFrom(err)
	}

	writeJSON(w, http.StatusOK, SuccessResponse{
		Status:  "success",
		Message: "User added to organisation successfully",
	})

	return nil
}
