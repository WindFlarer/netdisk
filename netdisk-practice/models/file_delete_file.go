package models

type FileDeleteFileRequest struct {
	Path string `json:"path" form:"path"`
}
