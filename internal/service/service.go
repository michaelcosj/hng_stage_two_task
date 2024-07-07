package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/michaelcosj/hng-task-two/internal/db"
	database "github.com/michaelcosj/hng-task-two/internal/db"
)

type RegisterParams struct {
	Email     string
	FirstName string
	LastName  string
	Password  string
	Phone     string
}

type LoginParams struct {
	Email    string
	Password string
}

type CreateOrgParam struct {
	Name        string
	Description string
}

type service struct {
	repo database.RepoQuerier
}

type Service interface {
	Register(ctx context.Context, param RegisterParams) (*AuthData, error)
	Login(ctx context.Context, param LoginParams) (*AuthData, error)
	GetUser(ctx context.Context, authUserId uuid.UUID, userId uuid.UUID) (*UserData, error)
	GetUserOrganisations(ctx context.Context, userId uuid.UUID) (*OrgsData, error)
	GetUserOrganisationById(ctx context.Context, userId uuid.UUID, orgId uuid.UUID) (*OrgData, error)
	CreateOrganisation(ctx context.Context, userId uuid.UUID, param CreateOrgParam) (*OrgData, error)
	AddUserToOrganisation(ctx context.Context, orgId uuid.UUID, userId uuid.UUID) error
}

func New(repo db.RepoQuerier) Service {

	return &service{repo}
}
