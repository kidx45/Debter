package domain

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashedPassword"`
	FullName       string `json:"fullName"`
	Email          string `json:"email"`
}