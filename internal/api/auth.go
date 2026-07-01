// Package api defines DTO's used by the server and client
package api

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RegistrationResponse struct {
	Token string `json:"token"`
}
