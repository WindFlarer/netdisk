package models

type FileDownloadRequest struct {
	FileName string `json:"fileName" form:"fileName"`
	Path     string `json:"path" form:"path"`
	DownPath string `json:"downPath" form:"downPath"`
}
