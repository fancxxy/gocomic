package models

import (
	"fmt"

	"github.com/fancxxy/gocomic/backend/config"
	// import mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// Init function create gocomic database and all tables
func Init() error {
	var err error
	database := config.Config.Database
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		database.User, database.Password, database.Host, database.Name))

	if err != nil {
		return err
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	// db.AutoMigrate(&User{}, &Comic{}, &Subscriber{}, &Chapter{}, &Picture{})
	// db.Model(&Subscriber{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	// db.Model(&Subscriber{}).AddForeignKey("comic_id", "comics(id)", "RESTRICT", "RESTRICT")
	// db.Model(&Chapter{}).AddForeignKey("comic_id", "comics(id)", "RESTRICT", "RESTRICT")
	// db.Model(&Picture{}).AddForeignKey("chapter_id", "chapters(id)", "RESTRICT", "RESTRICT")

	return err
}
