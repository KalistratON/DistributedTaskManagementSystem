package types

type UserMessage struct {
	Id       string `json:"id"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id     string `json:"id"`
	Status string `json:"status"`
	ErrMsg string `json:"err_msg"`
}
