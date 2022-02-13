package services

import (
	"github.com/streadway/amqp"
	"github.com/xStrato/video-encoder-golang/domain/entities"
)

type JobWorker struct {
	Job     entities.Job
	Message *amqp.Delivery
	Error   error
}

func NewJobWorker() *JobWorker {
	return &JobWorker{}
}

