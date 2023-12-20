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
	db.AutoMigrate(&model.Children{})
	db.AutoMigrate(&model.Eli{})
	db.AutoMigrate(&model.EliField{})
	db.AutoMigrate(&model.Work{})
	db.AutoMigrate(&model.WorkField{})
	db.AutoMigrate(&model.Vol{})
	db.AutoMigrate(&model.VolField{})
	db.AutoMigrate(&model.Learn{})
	db.AutoMigrate(&model.LearnField{})

	return db
}
