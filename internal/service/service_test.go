package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/michaelcosj/hng-task-two/internal/app"
	"github.com/michaelcosj/hng-task-two/internal/db"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}
}

func setupService(ctx context.Context) (*pgx.Conn, Service) {
	conn, err := pgx.Connect(ctx, os.Getenv("MOCK_PG_URI"))
	if err != nil {
		log.Fatalf("Error initialising database: %v", err)
	}

	testQueries := db.New(conn)
	testRepo := db.NewRepoQuerier(testQueries, conn)
	testService := New(testRepo)

	return conn, testService
}

// Unit test 1
func TestRegisterService(t *testing.T) {
	ctx := context.Background()

	conn, testService := setupService(ctx)
	defer conn.Close(ctx)

	// test that user is created successfully
	// with correct details and a valid token
	// with correct user data is generated
	var test = struct {
		name  string
		input RegisterParams
		want  AuthData
	}{
		name: "Test User Registered Successfully and Token Details Is Valid",
		input: RegisterParams{
			Email:     "johnDoe@email.com",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "a password",
			Phone:     "01000000000",
		},
		want: AuthData{
			User: UserData{
				Email:     "johnDoe@email.com",
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "01000000000",
			},
		},
	}

	t.Run(test.name, func(t *testing.T) {
		data, err := testService.Register(context.Background(), test.input)
		if err != nil {
			t.Fatal(err)
		}

		token, err := app.VerifyToken(data.Token)
		if err != nil {
			t.Fatal(err)
		}

		mapClaims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			t.Fatal(fmt.Errorf("error getting token claims"))
		}

		userIdStr, ok := mapClaims["id"].(string)
		if !ok {
			t.Fatal(fmt.Errorf("error getting token data: %v", mapClaims))
		}

		userId, err := uuid.Parse(userIdStr)
		if err != nil {
			t.Fatal(err)
		}

		if data.User.Id != userId.String() {
			t.Errorf("token does not contain correct user details: want %s, got %s", data.User.Id, userId.String())
		}

		if data.User.FirstName != test.input.FirstName {
			t.Errorf("user first name invalid: want %s, got %s", data.User.FirstName, test.input.FirstName)
		}

		if data.User.LastName != test.input.LastName {
			t.Errorf("user last name invalid: want %s, got %s", data.User.LastName, test.input.LastName)
		}

		if data.User.Phone != test.input.Phone {
			t.Errorf("user phone invalid: want %s, got %s", data.User.Phone, test.input.Phone)
		}

		if data.User.Email != test.input.Email {
			t.Errorf("user email invalid: want %s, got %s", data.User.Email, test.input.Email)
		}
	})
}

func TestGetOrganisation(t *testing.T) {
	ctx := context.Background()
	conn, testService := setupService(ctx)
	defer conn.Close(ctx)

	// seed seedData
	firstUserData := RegisterParams{
		Email:     "emailone@email.com",
		FirstName: "user",
		LastName:  "one",
		Password:  "password",
	}

	secondUserData := RegisterParams{
		Email:     "emailtwo@email.com",
		FirstName: "user",
		LastName:  "two",
		Password:  "password",
	}

	first, err := testService.Register(ctx, firstUserData)
	if err != nil {
		t.Fatal(err)
	}

	second, err := testService.Register(ctx, secondUserData)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Test user cannot get organisation they don't belong to", func(t *testing.T) {
		firstUserId, err := uuid.Parse(first.User.Id)
		if err != nil {
			t.Errorf("error parsing user id: %v", err)
		}

		firstUserOrgs, err := testService.GetUserOrganisations(ctx, firstUserId)
		if err != nil {
			t.Fatal(err)
		}

		orgId, err := uuid.Parse(firstUserOrgs.Orgs[0].Id)
		if err != nil {
			t.Errorf("error parsing org id: %v", err)
		}

		userId, err := uuid.Parse(second.User.Id)
		if err != nil {
			t.Errorf("error parsing user id: %v", err)
		}

		OrgData, err := testService.GetUserOrganisationById(ctx, userId, orgId)
		if err == nil {
			t.Errorf("user %s should not see org %s that belongs to user %s", userId, OrgData.Id, firstUserId)
		}
	})
}
