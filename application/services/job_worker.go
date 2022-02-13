package services

import (
	"encoding/json"
	"os"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"github.com/xStrato/video-encoder-golang/application/common"
	"github.com/xStrato/video-encoder-golang/domain/entities"
	"github.com/xStrato/video-encoder-golang/infrastructure/utils"
)

type JobWorker struct {
	input      chan amqp.Delivery
	out        chan common.JobWorkerResult
	JobService *JobService
}

func NewJobWorker(input chan amqp.Delivery, out chan common.JobWorkerResult, js *JobService) *JobWorker {
	return &JobWorker{
		input:      input,
		out:        out,
		JobService: js,
	}
}

func (j *JobWorker) Start(thread int) {
	for message := range j.input {
		if err := utils.IsJson(string(message.Body)); err != nil {
			j.out <- *common.NewJobWorkerResult(entities.Job{}, &message, err)
			continue
		}

		if err := json.Unmarshal(message.Body, &j.JobService.Job.Video); err != nil {
			j.out <- *common.NewJobWorkerResult(entities.Job{}, &message, err)
			continue
		}
		j.JobService.Job.Video.ID = uuid.NewV4().String()

		if err := j.JobService.Job.Video.IsValid(); err != nil {
			j.out <- *common.NewJobWorkerResult(entities.Job{}, &message, err)
			continue
		}

		if err := j.JobService.VideoService.InsertVideo(); err != nil {
			j.out <- *common.NewJobWorkerResult(entities.Job{}, &message, err)
			continue
		}

		j.JobService.Job.OutputBucket = os.Getenv("OUTPUT_BUCKET_NAME")
		j.JobService.Job.ID = uuid.NewV4().String()
		j.JobService.Job.Status = "STARTING"
		j.JobService.Job.CreatedAt = time.Now()

		if _, err := j.JobService.JobRepo.Insert(j.JobService.Job); err != nil {
			j.out <- *common.NewJobWorkerResult(entities.Job{}, &message, err)
			continue
		}

		if err := j.JobService.Start(); err != nil {
			j.out <- *common.NewJobWorkerResult(entities.Job{}, &message, err)
			continue
		}
		j.out <- *common.NewJobWorkerResult(*j.JobService.Job, &message, nil)
	}
}
