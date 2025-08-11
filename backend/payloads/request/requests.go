package request

type RegisterRequest struct{
	Email string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct{
	
}