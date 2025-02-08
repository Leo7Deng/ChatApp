package models

type User struct {
	ID              string `json:"id"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
}

type RegisterData struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
}
