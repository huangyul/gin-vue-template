package dto

type FileListQueryParam struct {
	FileName string `json:"file_name"`
	UserId   string `json:"user_id"`
	PageSize int    `json:"page_size"`
	Page     int    `json:"page"`
}
