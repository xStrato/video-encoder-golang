package services

import (
	"errors"
	"os"
	"strconv"

	"github.com/xStrato/video-encoder-golang/domain/entities"
	"github.com/xStrato/video-encoder-golang/infrastructure/database/repositories"
)

type JobService struct {
	Job          *entities.Job
	JobRepo      *repositories.JobRepository
	VideoService VideoService
}

func NewJobService(j *entities.Job, r *repositories.JobRepository, v VideoService) *JobService {
	return &JobService{
		Job:          j,
		JobRepo:      r,
		VideoService: v,
	}
}

func (j *JobService) Start() error {
	if err := j.changejobStatus("DOWNLOADING"); err != nil {
		return j.jobFail(err)
	}
	j.VideoService.Download(os.Getenv("INPUT_BUCKET_NAME"))

	if err := j.changejobStatus("FRAGMENTING"); err != nil {
		return j.jobFail(err)
	}
	j.VideoService.Fragment()

	if err := j.changejobStatus("ENCODING"); err != nil {
		return j.jobFail(err)
	}
	j.VideoService.Encode()

	//UPLOADING
	if err := j.Upload(); err != nil {
		return j.jobFail(err)
	}
	if err := j.changejobStatus("FINISHING"); err != nil {
		return j.jobFail(err)
	}
	if err := j.VideoService.Finish(); err != nil {
		return j.jobFail(err)
	}
	if err := j.changejobStatus("COMPLETED"); err != nil {
		return j.jobFail(err)
	}
	return nil
}

func (j *JobService) Upload() error {
	if err := j.changejobStatus("UPLOADING"); err != nil {
		return j.jobFail(err)
	}
	path := os.Getenv("LOCAL_STORAGE_PATH") + "/" + j.Job.VideoID
	uploadService := NewUploadService(os.Getenv("OUTPUT_BUCKET_NAME"), path)
	threads, err := strconv.Atoi(os.Getenv("CONCURRENT_UPLOADS"))
	if err != nil {
		return err
	}
	done := make(chan string)
	go uploadService.ProcessUpload(threads, done)
	if result := <-done; result != "upload completed" {
		return j.jobFail(errors.New(result))
	}
	return nil
}

func (j *JobService) changejobStatus(status string) error {
	j.Job.Status = status
	jj, err := j.JobRepo.Update(j.Job)
	if err != nil {
		return j.jobFail(err)
	}
	if job, ok := jj.(*entities.Job); ok {
		j.Job = job
	}
	return nil
}

func (j *JobService) jobFail(e error) error {

	j.Job.Status = "FAILED"
	j.Job.Error = e.Error()

	if _, err := j.JobRepo.Update(j.Job); err != nil {
		return err
	}
	return nil
}
