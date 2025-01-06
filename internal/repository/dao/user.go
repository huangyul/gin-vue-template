package dao

import "gorm.io/gorm"

type UserDao interface{}

type UserDaoImpl struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &UserDaoImpl{
		db: db,
	}
}

type User struct {
	Id       int64 `gorm:"primary_key;AUTO_INCREMENT"`
	Avatar   string
	Username string
	Nickname string
	Password string
}
