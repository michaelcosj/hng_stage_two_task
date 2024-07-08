package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/michaelcosj/hng-task-two/internal/db"
	"github.com/michaelcosj/hng-task-two/internal/server"
	"github.com/michaelcosj/hng-task-two/internal/service"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
}

var (
	registerRoute     = "/auth/register"
	loginRoute        = "/auth/login"
	getUserRoute      = "/api/users/%s"
	getAllOrgsRoute   = "/api/organisations"
	getSingleOrgRoute = "/api/organisations/%s"
	createOrgRoute    = "/api/organisations"
	addUserToOrg      = "/api/organisations/%s/users"
)

type authResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Token string `json:"accessToken"`
		User  struct {
			Id        string `json:"userId"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Email     string `json:"email"`
			Phone     string `json:"phone"`
		} `json:"user"`
	} `json:"data"`
}

func TestRegisterCases(t *testing.T) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, os.Getenv("MOCK_PG_URI"))
	if err != nil {
		t.Errorf("Error initialising database: %v", err)
	}
	defer conn.Close()

	querier := db.New(conn)
	repo := db.NewRepoQuerier(querier, conn)
	svc := service.New(repo)

	handler := server.RegisterRoutes(svc)

	tests := []struct {
		name        string
		FirstName   string `json:"firstName"`
		LastName    string `json:"lastName"`
		Email       string `json:"email"`
		Password    string `json:"password"`
		Phone       string `json:"phone"`
		shouldError bool
	}{
		{
			name:        "Test first name missing",
			FirstName:   "",
			LastName:    "Doe",
			Email:       "1@mail.com",
			Password:    "password",
			shouldError: true,
		},
		{
			name:        "Test last name missing",
			FirstName:   "John",
			LastName:    "",
			Email:       "2@mail.com",
			Password:    "password",
			shouldError: true,
		},
		{
			name:        "Test email missing",
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "",
			Password:    "password",
			shouldError: true,
		},
		{
			name:        "Test password missing",
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "3@mail.com",
			Password:    "",
			shouldError: true,
		}, {
			name:        "Test success",
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "4@mail.com",
			Password:    "password",
			shouldError: false,
		}, {
			name:        "Test duplicate email",
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "4@mail.com",
			Password:    "password",
			shouldError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonUser, _ := json.Marshal(test)
			req := httptest.NewRequest("POST", registerRoute, bytes.NewBuffer(jsonUser))
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req)

			var data authResp
			json.NewDecoder(resp.Body).Decode(&data)

			fmt.Println(data)

			if test.shouldError {
				if resp.Code != http.StatusUnprocessableEntity {
					t.Fatalf("expected 422 error got %d", resp.Code)
				}
				return
			}

			if resp.Code != http.StatusCreated {
				t.Fatalf("expected 201 got %d", resp.Code)
			}

			if data.Data.User.FirstName != test.FirstName {
				t.Fatalf("first name got %s expected %s", data.Data.User.FirstName, test.FirstName)
			}

			if data.Data.User.LastName != test.LastName {
				t.Fatalf("last name got %s expected %s", data.Data.User.LastName, test.LastName)
			}

			if data.Data.User.Email != test.Email {
				t.Fatalf("email got %s expected %s", data.Data.User.Email, test.Email)
			}

			if len(data.Data.Token) == 0 {
				t.Fatalf("invalid token")
			}
		})
	}
}

func TestLoginCases(t *testing.T) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, os.Getenv("MOCK_PG_URI"))
	if err != nil {
		t.Errorf("Error initialising database: %v", err)
	}
	defer conn.Close()

	querier := db.New(conn)
	repo := db.NewRepoQuerier(querier, conn)
	svc := service.New(repo)

	handler := server.RegisterRoutes(svc)

	tests := []struct {
		name        string
		Email       string `json:"email"`
		Password    string `json:"password"`
		shouldError bool
	}{
		{
			name:        "Test successful login",
			Email:       "4@mail.com",
			Password:    "password",
			shouldError: false,
		}, {
			name:        "Test failed login",
			Email:       "incorrect@mail.com",
			Password:    "password",
			shouldError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			jsonUser, _ := json.Marshal(test)
			req := httptest.NewRequest("POST", loginRoute, bytes.NewBuffer(jsonUser))
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req)

			var data authResp
			json.NewDecoder(resp.Body).Decode(&data)

			fmt.Println(data)

			if test.shouldError {
				if resp.Code != http.StatusUnauthorized {
					t.Fatalf("expected 401 error got %d", resp.Code)
				}
				return
			}

			if resp.Code != http.StatusOK {
				t.Fatalf("expected 201 got %d", resp.Code)
			}

			if data.Data.User.Email != test.Email {
				t.Fatalf("email got %s expected %s", data.Data.User.Email, test.Email)
			}

			if len(data.Data.Token) == 0 {
				t.Fatalf("invalid token")
			}
		})
	}
}

func TestOrgCases(t *testing.T) {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, os.Getenv("MOCK_PG_URI"))
	if err != nil {
		t.Errorf("Error initialising database: %v", err)
	}
	defer conn.Close()

	querier := db.New(conn)
	repo := db.NewRepoQuerier(querier, conn)
	svc := service.New(repo)

	handler := server.RegisterRoutes(svc)

	t.Run("Test Default Organisation Exists With Correct Name", func(t *testing.T) {
		// register new user to get their token
		jsonUser, _ := json.Marshal(&struct {
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
			Email     string `json:"email"`
			Password  string `json:"password"`
			Phone     string `json:"phone"`
		}{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "420@mail.com",
			Password:  "password",
		})

		regReq := httptest.NewRequest("POST", registerRoute, bytes.NewBuffer(jsonUser))
		regResp := httptest.NewRecorder()
		handler.ServeHTTP(regResp, regReq)

		var regData authResp
		json.NewDecoder(regResp.Body).Decode(&regData)

		req := httptest.NewRequest("GET", getAllOrgsRoute, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", regData.Data.Token))
		resp := httptest.NewRecorder()

		handler.ServeHTTP(resp, req)

		var data struct {
			Status  string `json:"Status"`
			Message string `json:"Message"`
			Data    struct {
				Orgs []struct {
					OrgId       string `json:"OrgId"`
					Name        string `json:"Name"`
					Description string `json:"Description"`
				} `json:"organisations"`
			} `json:"data"`
		}

		json.NewDecoder(resp.Body).Decode(&data)

		if resp.Code != http.StatusOK {
			t.Fatalf("expected 201 got %d", resp.Code)
		}

		expectedOrgName := fmt.Sprintf("%s's Organisation", regData.Data.User.FirstName)
		if data.Data.Orgs[0].Name != expectedOrgName {
			t.Fatalf("email expected %s got %s", expectedOrgName, data.Data.Orgs[0].Name)
		}
	})
}
