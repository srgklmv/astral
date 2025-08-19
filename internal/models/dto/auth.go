package dto

type (
	RegisterRequest struct {
		Token    string `json:"token"`
		Login    string `json:"login"`
		Password string `json:"pswd"`
	}
	RegisterResponse struct {
		Login string `json:"login"`
	}
)

type AuthResponse struct {
	Token string `json:"token"`
}

type LogoutResponse map[string]bool
