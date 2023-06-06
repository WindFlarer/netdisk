package models

type FileDeleteDirRequest struct {
	FileName string `json:"fileName" form:"fileName"`
	Path     string `json:"path" form:"path"`
}
