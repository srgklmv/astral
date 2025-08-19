package dto

type RegisterResponse struct {
	Login string `json:"login"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type LogoutResponse map[string]bool
