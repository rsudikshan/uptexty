package models

type ResponseModel struct{
	Success bool `json:"status"`
	Message string `json:"message"`
	Body any `json:"body"`
}