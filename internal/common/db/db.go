package db

import (
	"log"

	"github.com/jcalabing/hrmis-api/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.UserField{})
	db.AutoMigrate(&model.Edu{})
	db.AutoMigrate(&model.EduField{})

	return db
}
