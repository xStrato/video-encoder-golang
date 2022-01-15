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
	OutputBucket  string `json:"output_bucket" valid:"required"`
	Status        string `json:"status" valid:"required"`
	Video         *Video `json:"video" valid:"-"`
	VideoID       string `json:"-" valid:"-" gorm:"column:video_id;type:uuid;notnull"`
	Error         string `json:"error" valid:"-"`
}

func NewJob(output, status string, v *Video) *Job {
	return &Job{
		Entity:       *common.NewEntity(),
		OutputBucket: output,
		Status:       status,
		Video:        v,
	}
}

func (j *Job) IsEntity() bool {
	return true
}
