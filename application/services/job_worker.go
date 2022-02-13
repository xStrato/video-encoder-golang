package services

import (
	"encoding/json"
	"os"
	"sync"
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
	jobService *JobService
	mutex      sync.Mutex
}

func NewJobWorker(input chan amqp.Delivery, out chan common.JobWorkerResult, js *JobService) *JobWorker {
	return &JobWorker{
		input:      input,
		out:        out,
		jobService: js,
		mutex:      sync.Mutex{},
	}
}

func (j *JobWorker) Start(thread int) {
	for message := range j.input {
		if err := utils.IsJson(string(message.Body)); err != nil {
			j.out <- *common.NewJobWorkerResult(entities.Job{}, &message, err)
			continue
		}

		if err := json.Unmarshal(message.Body, &j.jobService.Job.Video); err != nil {
			j.out <- *common.NewJobWorkerResult(entities.Job{}, &message, err)
			continue
		}
		j.jobService.Job.Video.ID = uuid.NewV4().String()

		if err := j.jobService.Job.Video.IsValid(); err != nil {
			j.out <- *common.NewJobWorkerResult(entities.Job{}, &message, err)
			continue
		}

		if err := j.jobService.VideoService.InsertVideo(); err != nil {
			j.out <- *common.NewJobWorkerResult(entities.Job{}, &message, err)
			continue
		}

		j.jobService.Job.OutputBucket = os.Getenv("OUTPUT_BUCKET_NAME")
		j.jobService.Job.ID = uuid.NewV4().String()
		j.jobService.Job.Status = "STARTING"
		j.jobService.Job.CreatedAt = time.Now()

		if _, err := j.jobService.JobRepo.Insert(j.jobService.Job); err != nil {
			j.out <- *common.NewJobWorkerResult(entities.Job{}, &message, err)
			continue
		}

		if err := j.jobService.Start(); err != nil {
			j.out <- *common.NewJobWorkerResult(entities.Job{}, &message, err)
			continue
		}
		j.out <- *common.NewJobWorkerResult(*j.jobService.Job, &message, nil)
	}
}
