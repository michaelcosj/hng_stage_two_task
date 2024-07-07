package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/michaelcosj/hng-task-two/internal/app"
	"github.com/michaelcosj/hng-task-two/internal/db"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func (s *service) Register(ctx context.Context, param RegisterParams) (*AuthData, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	// use a transaction so a user cannot be created if
	// creating and attaching it to it's default organisation fails
	tx, err := s.repo.GetDB().Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot create database transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	qTx := s.repo.WithTx(tx)

	user, err := qTx.UserInsert(ctx, db.UserInsertParams{
		Email:     param.Email,
		FirstName: param.FirstName,
		LastName:  param.LastName,
		Password:  string(passwordHash),
		Phone:     pgtype.Text{String: param.Phone, Valid: len(param.Phone) == 11},
	})

	if err != nil {
		// if a unique violation error occurs, it means a user with this email already exists
		// this is a user request error
		if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == pgerrcode.UniqueViolation {
			return nil, app.ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("error in user registration service: %w", err)
	}

	// create a default organisation for the user
	// using the user's name to generate the organisation name
	fNameCapitalised := cases.Title(language.English, cases.Compact).String(param.FirstName)
	org, err := qTx.OrgInsert(ctx, db.OrgInsertParams{
		Name:        fmt.Sprintf("%s's Organisation", fNameCapitalised),
		Description: pgtype.Text{Valid: false},
	})

	if err != nil {
		return nil, fmt.Errorf("error in user registration service: %w", err)
	}

	// add the user to the default organisation
	if err = qTx.UserAddOrg(ctx, db.UserAddOrgParams{
		UserID: user.ID,
		OrgID:  org.ID,
	}); err != nil {
		return nil, fmt.Errorf("error in user registration service: %w", err)
	}

	// commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// create jwt token
	token, err := app.CreateToken(user.ID.String())
	if err != nil {
		return nil, fmt.Errorf("error in user registration service: %w", err)
	}

	return &AuthData{
		Token: token,
		User: UserData{
			Id:        user.ID.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Phone:     user.Phone.String,
		},
	}, nil
}

func (s *service) Login(ctx context.Context, param LoginParams) (*AuthData, error) {
	user, err := s.repo.UserWhereEmail(ctx, param.Email)
	if err != nil {
		return nil, app.ApiErrorFrom(fmt.Errorf("error retrieving user from db: %w", app.ErrAuthenticationFailed))
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(param.Password))
	if err != nil {
		return nil, app.ApiErrorFrom(fmt.Errorf("error comparing user password with hash: %w", app.ErrAuthenticationFailed))
	}

	// create jwt token
	token, err := app.CreateToken(user.ID.String())
	if err != nil {
		return nil, fmt.Errorf("error in user registration service: %w", err)
	}

	return &AuthData{
		Token: token,
		User: UserData{
			Id:        user.ID.String(),
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Phone:     user.Phone.String,
		},
	}, nil
}
