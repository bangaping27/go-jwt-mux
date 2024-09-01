package models

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/go_jwt_mux"))
	if err != nil {
		fmt.Println("Gagal terhubung ke database: ", err)
	}

	db.AutoMigrate(&User{})
	DB = db
}
