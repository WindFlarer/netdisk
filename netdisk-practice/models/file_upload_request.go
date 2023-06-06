package models

type FileUploadRequest struct {
	FileName string `json:"fileName" form:"fileName"`
	Path     string `json:"path" form:"path"`
}
