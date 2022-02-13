package cmd

import (
	"github.com/streadway/amqp"
	"github.com/xStrato/video-encoder-golang/application/common"
	"github.com/xStrato/video-encoder-golang/application/services"
	"github.com/xStrato/video-encoder-golang/infrastructure/database/contexts"
	"github.com/xStrato/video-encoder-golang/infrastructure/queue"
	"gorm.io/gorm"
)

func StartServer(conn *gorm.DB) {
	input := make(chan amqp.Delivery)
	output := make(chan common.JobWorkerResult)

	rabbitMQ := queue.NewRabbitMQ()
	channel := rabbitMQ.Connect()
	defer channel.Close()

	rabbitMQ.Consume(input)
	dbContext := contexts.NewEncoderContext(conn)
	jobManager := services.NewJobManager(*dbContext, input, output, rabbitMQ)
	jobManager.Start()
}
