package models

import "gorm.io/gorm"

type UserBasic struct {
	ID       uint   `gorm:"primaryKey; column:id; type:int;" json:"id"`
	UserName string `gorm:"column:user_name; type:varchar(255);" json:"userName"`
	Password string `gorm:"column:password; type:varchar(255);" json:"password"`
	Email    string `gorm:"column:email; type:varchar(255);" json:"email"`
	Phone    string `gorm:"column:phone; type:varchar(11);" json:"phone"`
	gorm.Model
}

func (u *UserBasic) TableName() string {
	return "user_basic"
}
