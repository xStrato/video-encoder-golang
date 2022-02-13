package common

import (
	"github.com/streadway/amqp"
	"github.com/xStrato/video-encoder-golang/domain/entities"
)

type JobWorkerResult struct {
	Job     entities.Job
	Message *amqp.Delivery
	Error   error
}

func NewJobWorkerResult(j entities.Job, m *amqp.Delivery, e error) *JobWorkerResult {
	return &JobWorkerResult{
		Job:     j,
		Message: m,
		Error:   e,
	}
}
