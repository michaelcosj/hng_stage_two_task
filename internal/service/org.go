package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/michaelcosj/hng-task-two/internal/app"
	"github.com/michaelcosj/hng-task-two/internal/db"
)

func (s *service) GetUserOrganisations(ctx context.Context, userId uuid.UUID) (*OrgsData, error) {
	orgs, err := s.repo.OrgAllWhereUser(ctx, userId)
	if err != nil {
		return nil, app.ApiErrorFrom(fmt.Errorf("error finding user orgs: %v", app.ErrClientError))
	}

	resp := &OrgsData{}
	for _, org := range orgs {
		resp.Orgs = append(resp.Orgs, OrgData{
			Id:          org.ID.String(),
			Name:        org.Name,
			Description: org.Description.String,
		})
	}

	return resp, nil
}

func (s *service) GetUserOrganisationById(ctx context.Context, userId uuid.UUID, orgId uuid.UUID) (*OrgData, error) {
	org, err := s.repo.OrgWhereUser(ctx, db.OrgWhereUserParams{
		UserID: userId,
		OrgID:  orgId,
	})

	if err != nil {
		return nil, app.ApiErrorFrom(fmt.Errorf("error retrieving user organisation from db: %v", app.ErrOrgNotFound))
	}

	return &OrgData{
		Id:          org.ID.String(),
		Name:        org.Name,
		Description: org.Description.String,
	}, nil
}

func (s *service) CreateOrganisation(ctx context.Context, userId uuid.UUID, param CreateOrgParam) (*OrgData, error) {
	// create organisation in a transaction
	// for the same reason as in register service
	tx, err := s.repo.GetDB().Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot create database transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	qTx := s.repo.WithTx(tx)

	org, err := s.repo.OrgInsert(ctx, db.OrgInsertParams{
		Name:        param.Name,
		Description: pgtype.Text{String: param.Description, Valid: len(param.Description) != 0},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating organisation: %v", err)
	}

	if err = qTx.UserAddOrg(ctx, db.UserAddOrgParams{
		UserID: userId,
		OrgID:  org.ID,
	}); err != nil {
		return nil, app.ApiErrorFrom(fmt.Errorf("error adding user to organisation: %v", app.ErrUserNotFound))
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &OrgData{
		Id:          org.ID.String(),
		Name:        org.Name,
		Description: org.Description.String,
	}, nil
}

func (s *service) AddUserToOrganisation(ctx context.Context, orgId uuid.UUID, userId uuid.UUID) error {
	// check if user exists
	if _, err := s.repo.UserWhereId(ctx, userId); err != nil {
		return app.ErrUserNotFound
	}

	err := s.repo.UserAddOrg(ctx, db.UserAddOrgParams{
		UserID: userId,
		OrgID:  orgId,
	})

	if err != nil {
		return app.ErrOrgNotFound
	}

	return nil
}
