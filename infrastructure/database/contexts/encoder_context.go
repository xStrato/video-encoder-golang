package contexts

import (
	"log"

	"github.com/xStrato/video-encoder-golang/domain/entities"
	"gorm.io/gorm"
)

type EncoderContext struct {
	db *gorm.DB
}

func NewEncoderContext(db *gorm.DB) *EncoderContext {
	return &EncoderContext{db}
}

func (v *EncoderContext) RunMigrations() {
	if err := v.db.AutoMigrate(&entities.Video{}); err != nil {
		log.Fatalln("Cannot run migrations for Video entity: ", err.Error())
	}
	if err := v.db.AutoMigrate(&entities.Job{}); err != nil {
		log.Fatalln("Cannot run migrations for Job entity: ", err.Error())
	}
}

func (v *EncoderContext) GetDBConnection() *gorm.DB {
	return v.db
}
