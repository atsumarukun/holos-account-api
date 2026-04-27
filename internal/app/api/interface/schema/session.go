package schema

import "github.com/google/uuid"

type CreateSessionRequest struct {
	AccountName string `json:"account_name"`
	Password    string `json:"password"`
}

type SessionResponse struct {
	Token string `json:"token"`
}

type VerifySessionResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
