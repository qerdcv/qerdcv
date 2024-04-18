package handlers

type RegisterResponse struct {
	Message string `json:"message"`
}

type AuthorizeUserResponse struct {
	SessionToken string `json:"session_token"`
}
