package response

type LoginResponse struct{
	Jwt string `json:"jwt"`
	Refresh string `json:"refresh"`
}