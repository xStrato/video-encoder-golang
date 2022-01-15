package common

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Entity struct {
	ID        string    `valid:"uuid"`
	CreatedAt time.Time `valid:"-"`
	UpdatedAt time.Time `valid:"-"`
}

func NewEntity() *Entity {
	return &Entity{
		ID:        uuid.NewV4().String(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (e *Entity) IsValid() error {
	if valid, err := govalidator.ValidateStruct(e); !valid {
		return err
	}
	return nil
}
