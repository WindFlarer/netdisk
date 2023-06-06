package models

import "gorm.io/gorm"

type FileBasic struct {
	ID       uint   `gorm:"primaryKey; column:id; type:int;" json:"id"`
	COSPath  string `gorm:"column:cos_path; type:varchar(255);" json:"COSPath"`
	FileName string `gorm:"column:file_name; type:varchar(255);" json:"fileName"`
	UserName string `gorm:"column:user_name; type:varchar(255);" json:"userName"`
	FileSize int64  `gorm:"column:file_size; type:bigint;" json:"fileSize"`
	Path     string `gorm:"column:path; type:varchar(255);" json:"path"`
	IsDir    bool   `gorm:"column:is_dir; type:bool;" json:"isDir"`
	gorm.Model
}

func (f *FileBasic) TableName() string {
	return "file_basic"
}
