package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDB() {

	dsn := "root:@tcp(127.0.0.1:3306)/minihifpt?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := "root:@tcp(127.0.0.1:3306)/testminihifpt?charset=utf8mb4&parseTime=True&loc=Local"

	// Kết nối tới DB
	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Không thể kết nối cơ sở dữ liệu: ", err)
	}

	DB = connection
	log.Println("Kết nối cơ sở dữ liệu thành công!")
}
