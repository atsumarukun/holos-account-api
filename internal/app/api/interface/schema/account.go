package schema

type CreateAccountRequest struct {
	Name            string `json:"name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type AccountResponse struct {
	Name string `json:"name"`
}
