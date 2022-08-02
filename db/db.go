package db

import (
	"awesomeProject/domain"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var err error

func InitializeDb() {
	db, err = gorm.Open(postgres.Open("host=localhost user=root password=root dbname=DeckMS port=5455 sslmode=disable TimeZone=Asia/Yerevan"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		fmt.Println(err.Error())
		panic("Error connecting to db")
	}

}

func MigrateTables() {
	db.AutoMigrate(&domain.Deck{}, &domain.Card{})
}

func GetDbInstance() *gorm.DB {
	return db
}
