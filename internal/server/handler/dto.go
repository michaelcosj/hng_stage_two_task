package handler

import (
	"net/mail"
)

// Validation in this application is basic
// only to conform to the task specifications
type RegisterUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Phone     string `json:"phone"`
}

func (req *RegisterUserRequest) Validate() map[string]string {
	problems := make(map[string]string)

	if len(req.Email) == 0 {
		problems["email"] = "email must be provided"
	}

	if !isValidEmail(req.Email) {
		problems["email"] = "email is invalid"
	}

	if len(req.FirstName) == 0 {
		problems["firstName"] = "first name must be provided"
	}

	if len(req.LastName) == 0 {
		problems["lastName"] = "last name must be provided"
	}

	if len(req.Password) == 0 {
		problems["password"] = "password must be provided"
	}

	return problems
}

type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req *LoginUserRequest) Validate() map[string]string {
	problems := make(map[string]string)

	if len(req.Email) == 0 {
		problems["email"] = "email must be provided"
	}

	if !isValidEmail(req.Email) {
		problems["email"] = "email is invalid"
	}

	if len(req.Password) == 0 {
		problems["password"] = "password must be provided"
	}

	return problems
}

type CreateOrgRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (req *CreateOrgRequest) Validate() map[string]string {
	problems := make(map[string]string)

	if len(req.Name) == 0 {
		problems["name"] = "name must be provided"
	}

	return problems
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
