package service

type UserData struct {
	Id        string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type AuthData struct {
	Token string   `json:"accessToken"`
	User  UserData `json:"user"`
}

type OrgData struct {
	Id          string `json:"orgId"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type OrgsData struct {
	Orgs []OrgData `json:"organisations"`
}
