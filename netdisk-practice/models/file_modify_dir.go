package models

type FileModifyDirRequest struct {
	NewDirName string `json:"newDirName" form:"newDirName"`
	OldDirName string `json:"oldDirName" form:"oldDirName"`
	Path       string `json:"path" form:"path"`
}
