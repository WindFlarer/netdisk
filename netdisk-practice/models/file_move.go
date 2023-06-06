package models

type FileMoveRequest struct {
	OldPath string `json:"oldPath" form:"oldPath"`
	NewPath string `json:"newPath" form:"newPath"`
}
