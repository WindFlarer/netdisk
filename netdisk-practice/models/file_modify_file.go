package models

type FileModifyFileRequest struct {
	NewFileName string `json:"newFileName" form:"newFileName"`
	Path        string `json:"path" form:"path"`
}
