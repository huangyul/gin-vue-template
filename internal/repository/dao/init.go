package dao

import "gorm.io/gorm"

func InitTable(db *gorm.DB) {
	err := db.AutoMigrate(&User{}, &File{})
	if err != nil {
		panic(err)
	}
}
