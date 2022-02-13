package services

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/streadway/amqp"
	"github.com/xStrato/video-encoder-golang/application/common"
	"github.com/xStrato/video-encoder-golang/domain/entities"
	"github.com/xStrato/video-encoder-golang/infrastructure/database/contexts"
	"github.com/xStrato/video-encoder-golang/infrastructure/database/repositories"
	"github.com/xStrato/video-encoder-golang/infrastructure/queue"
)

type JobManager struct {
	dbContext contexts.EncoderContext
	job       entities.Job
	inputs    chan amqp.Delivery
	output    chan common.JobWorkerResult
	rabbitMQ  *queue.RabbitMQ
}

func NewJobManager(c contexts.EncoderContext, i chan amqp.Delivery, o chan common.JobWorkerResult, r *queue.RabbitMQ) *JobManager {
	return &JobManager{
		dbContext: c,
		job:       entities.Job{},
		inputs:    i,
		output:    o,
		rabbitMQ:  r,
	}
}
func (j *JobManager) Start() error {
	videoRepo := repositories.NewVideoRepository(&j.dbContext)
	videoService := NewVideoService(videoRepo)

	jobRepo := repositories.NewJobRepository(&j.dbContext)
	jobService := NewJobService(&j.job, jobRepo, *videoService)

	workers, err := strconv.Atoi(os.Getenv("CONCURRENT_WORKERS"))
	if err != nil {
		return err
	}

	jobWorker := NewJobWorker(j.inputs, j.output, jobService)
	for thread := 0; thread < workers; thread++ {
		go jobWorker.Start(thread)
	}

	for jobResult := range j.output {
		if jobResult.Error != nil {
			err = j.parseErrors(&jobResult)
		} else {
			err = j.notifySuccess(&jobResult)
		}

		if err != nil {
			jobResult.Message.Reject(false)
		}
	}
	return nil
}

func (j *JobManager) notifySuccess(jr *common.JobWorkerResult) error {
	json, err := json.Marshal(jr.Job)
	if err != nil {
		return err
	}
	if err := j.notify(json); err != nil {
		return err
	}
	if err := jr.Message.Ack(false); err != nil {
		return err
	}
	return nil
}

func (j *JobManager) parseErrors(jr *common.JobWorkerResult) error {
	if jr.Job.ID != "" {
		log.Printf("MessageID: %v. Error during the job: %v with video: %v. Error: %v", jr.Message.DeliveryTag, jr.Job.ID, jr.Job.Video.ID, jr.Error.Error())
	} else {
		log.Printf("MessageID: %v. Error parsing message: %v", jr.Message.DeliveryTag, jr.Error)
	}

	errMsg := common.NewJobNotification(string(jr.Message.Body), jr.Error.Error())
	json, err := json.Marshal(errMsg)
	if err != nil {
		return err
	}
	if err := j.notify(json); err != nil {
		return err
	}
	if err := jr.Message.Reject(false); err != nil {
		return err
	}
	return nil
}

func (j *JobManager) notify(json []byte) error {
	err := j.rabbitMQ.Notify(
		string(json),
		"application/json",
		os.Getenv("RABBITMQ_NOTIFICATION_EX"),
		os.Getenv("RABBITMQ_NOTIFICATION_ROUTING_KEY"),
	)
	if err != nil {
		return err
	}
	return nil
}
