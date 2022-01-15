package entities

import (
	"github.com/asaskevich/govalidator"
	"github.com/xStrato/video-encoder-golang/domain/common"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Job struct {
	common.Entity `valid:"required"`
	OutputBucket  string `valid:"required"`
	Status        string `valid:"required"`
	Video         *Video `valid:"-"`
	VideoID       string `valid:"-"`
	Error         string `valid:"-"`
}

func NewJob(output, status string, v *Video) *Job {
	return &Job{
		Entity:       *common.NewEntity(),
		OutputBucket: output,
		Status:       status,
		Video:        v,
	}
}
