package models

type User struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// struct for authentification
type AuthDto struct {
	Email    string
	Password string
}

// struct for requests that dont require password
type UserResponse struct {
	Id    uint64 `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
