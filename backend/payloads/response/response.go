package response

type LoginResponse struct{
	Jwt string `json:"jwt"`
	Refresh string `json:"refresh"`
}

type GetFilesResponse struct {
	ID int `json:"id"`
	Filename string `json:"filename"`
	UploadedAt string `json:"uploaded_at"`
}

type GetRowsResponse struct{
	Id int `json:"id"`
	Position float64 `json:"position"`
	InputText string `json:"input_text"`
}