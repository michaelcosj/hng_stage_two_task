package mock

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/michaelcosj/hng-task-two/internal/db"
)

type MockRepo struct {
	userStore    map[string]db.User
	orgStore     map[string]db.Organisation
	userOrgStore map[string][]db.UserOrganisation

	db MockDB
}

func NewMockRepo() *MockRepo {
	return &MockRepo{db: MockDB{}}
}

func (r *MockRepo) WithTx(tx pgx.Tx) db.RepoQuerier {
	return r
}

func (r *MockRepo) GetDB() db.Db {
	return r.db
}

func (r *MockRepo) FindUserInOrgs(ctx context.Context, arg db.FindUserInOrgsParams) (db.User, error) {
	user, err := r.UserWhereId(ctx, arg.FindUser)
	if err != nil {
		return db.User{}, err
	}

	for _, auth_user_org := range r.userOrgStore[arg.AuthUser.String()] {
		auth_org, _ := r.OrganisationWhereId(ctx, auth_user_org.OrgID)
		for _, find_user_org := range r.userOrgStore[arg.FindUser.String()] {
			find_org, _ := r.OrganisationWhereId(ctx, find_user_org.OrgID)
			if auth_org == find_org {
				return user, nil
			}
		}
	}

	return db.User{}, fmt.Errorf("not found")
}

func (r *MockRepo) OrgAllWhereUser(ctx context.Context, userID uuid.UUID) ([]db.Organisation, error) {
	var orgs []db.Organisation
	for _, user_org := range r.userOrgStore[userID.String()] {
		org, err := r.OrganisationWhereId(ctx, user_org.OrgID)
		if err == nil {
			orgs = append(orgs, org)
		}
	}

	return orgs, nil
}

func (r *MockRepo) OrgInsert(ctx context.Context, arg db.OrgInsertParams) (db.Organisation, error) {
	org := db.Organisation{
		ID:          uuid.New(),
		Name:        arg.Name,
		Description: arg.Description,
	}

	r.orgStore[org.ID.String()] = org

	return org, nil
}

func (r *MockRepo) OrgWhereUser(ctx context.Context, arg db.OrgWhereUserParams) (db.Organisation, error) {
	for _, user_org := range r.userOrgStore[arg.UserID.String()] {
		if user_org.OrgID.String() == arg.OrgID.String() {
			return r.OrganisationWhereId(ctx, user_org.OrgID)
		}
	}

	return db.Organisation{}, fmt.Errorf("not found")
}

func (r *MockRepo) OrganisationWhereId(ctx context.Context, id uuid.UUID) (db.Organisation, error) {
	org, ok := r.orgStore[id.String()]
	if !ok {
		return db.Organisation{}, fmt.Errorf("not found")
	}

	return org, nil
}

func (r *MockRepo) UserAddOrg(ctx context.Context, arg db.UserAddOrgParams) error {
	user_org := db.UserOrganisation{
		UserID: arg.UserID,
		OrgID:  arg.OrgID,
	}

	r.userOrgStore[user_org.UserID.String()] = append(r.userOrgStore[user_org.UserID.String()], user_org)
	return nil
}

func (r *MockRepo) UserInsert(ctx context.Context, arg db.UserInsertParams) (db.User, error) {
	user := db.User{
		ID:        uuid.New(),
		Email:     arg.Email,
		FirstName: arg.FirstName,
		LastName:  arg.LastName,
		Password:  arg.Password,
		Phone:     arg.Phone,
	}

	r.userStore[user.ID.String()] = user
	return user, nil
}

func (r *MockRepo) UserWhereEmail(ctx context.Context, email string) (db.User, error) {
	for _, user := range r.userStore {
		if user.Email == email {
			return user, nil
		}
	}

	return db.User{}, fmt.Errorf("not found")
}

func (r *MockRepo) UserWhereId(ctx context.Context, id uuid.UUID) (db.User, error) {
	user, ok := r.userStore[id.String()]
	if !ok {
		return db.User{}, fmt.Errorf("not found")
	}

	return user, nil
}
