package schema

import "github.com/google/uuid"

type LoginRequest struct {
	AccountName string `json:"account_name"`
	Password    string `json:"password"`
}

type SessionResponse struct {
	Token string `json:"token"`
}

type AauthorizationResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
