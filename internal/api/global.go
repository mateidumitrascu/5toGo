package api

type ErrorResponse struct {
	ErrorMessage string `json:"error"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
