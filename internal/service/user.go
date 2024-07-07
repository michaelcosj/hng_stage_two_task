package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/michaelcosj/hng-task-two/internal/app"
	"github.com/michaelcosj/hng-task-two/internal/db"
)

func (s *service) GetUser(ctx context.Context, authUserId uuid.UUID, userId uuid.UUID) (*UserData, error) {
	var err error
	var user db.User
	if authUserId.String() == userId.String() {
		// if userid is the authenticated user's id
		// we return the authenticated user's data
		user, err = s.repo.UserWhereId(ctx, authUserId)
	} else {
		// if not, we find the user with the id that
		// belongs to an organisation that the authenticated user
		// also belongs to
		user, err = s.repo.FindUserInOrgs(ctx, db.FindUserInOrgsParams{
			AuthUser: authUserId,
			FindUser: userId,
		})
	}

	if err != nil {
		return nil, app.ApiErrorFrom(fmt.Errorf("error getting user from db: %w", app.ErrUserNotFound))
	}

	return &UserData{
		Id:        user.ID.String(),
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Phone:     user.Phone.String,
	}, nil
}
