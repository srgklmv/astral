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

type (
	AuthRequest struct {
		Login    string `json:"login"`
		Password string `json:"pswd"`
	}
	AuthResponse struct {
		Token string `json:"token"`
	}
)

type (
	LogoutRequest struct {
		Token string `params:"token"`
	}
	LogoutResponse map[string]bool
)
