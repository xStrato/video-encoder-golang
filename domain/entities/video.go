package entities

import (
	"github.com/asaskevich/govalidator"
	"github.com/xStrato/video-encoder-golang/domain/common"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Video struct {
	common.Entity `valid:"required"`
	ResourceID    string `json:"resource_id" valid:"uuid"`
	FilePath      string `json:"file_path" valid:"required"`
	Jobs          []*Job `json:"-" valid:"-" gorm:"foreignKey:VideoID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func NewVideo(rID, file string) *Video {
	return &Video{
		Entity:     *common.NewEntity(),
		ResourceID: rID,
		FilePath:   file,
	}
}

func (v *Video) IsEntity() bool {
	return true
}
