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
	ResourceID    string `valid:"uuid"`
	FilePath      string `valid:"required"`
}

func NewVideo(rID, file string) *Video {
	return &Video{
		Entity:     *common.NewEntity(),
		ResourceID: rID,
		FilePath:   file,
	}
}
