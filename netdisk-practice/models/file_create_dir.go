package models

type FileCreateDir struct {
	Path     string `json:"path" form:"path"`
	FileName string `json:"fileName" form:"fileName"`
}
