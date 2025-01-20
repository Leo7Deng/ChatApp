package models

type RegisterData struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}