package api

// AuthRegisterRequest POST /v1/auth/register/
type AuthRegisterRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

// AuthLoginRequest POST /v1/auth/login/
type AuthLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthLoginResponse /v1/auth/login/
type AuthLoginResponse tokenResponse

// AuthPasswdRequest POST /v1/auth/passwd/
type AuthPasswdRequest struct {
	Username    string `json:"username,omitempty"`
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

// AuthRegenerateRequest POST /v1/auth/tokens/
type AuthRegenerateRequest struct {
	Name string `json:"username,omitempty"`
	All  bool   `json:"all,omitempty"`
}

// AuthRegenerateResponse /v1/auth/tokens/
type AuthRegenerateResponse tokenResponse

type tokenResponse struct {
	Token string `json:"token"`
}
