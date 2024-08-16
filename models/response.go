package models

type AuthResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type Response struct {
	Message string `json:"message"`
}
