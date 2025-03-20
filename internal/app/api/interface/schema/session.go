package schema

type LoginRequest struct {
	AccountName string `json:"account_name"`
	Password    string `json:"password"`
}

type SessionResponse struct {
	Token string `json:"token"`
}
