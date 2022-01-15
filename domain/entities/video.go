package entities

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Video struct {
	ID         string    `valid:"uuid"`
	ResourceID string    `valid:"uuid"`
	FilePath   string    `valid:"required"`
	CreatedAt  time.Time `valid:"-"`
}

func NewVideo(rID, file string) *Video {
	return &Video{
		ID:         uuid.NewV4().String(),
		ResourceID: rID,
		FilePath:   file,
		CreatedAt:  time.Now(),
	}
}

func (v *Video) IsValid() error {
	if valid, err := govalidator.ValidateStruct(v); !valid {
		return err
	}
	return nil
}
