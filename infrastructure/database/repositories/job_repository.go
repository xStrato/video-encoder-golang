package repositories

import (
	"fmt"

	"github.com/xStrato/video-encoder-golang/domain/entities"
	"github.com/xStrato/video-encoder-golang/domain/interfaces"
	"github.com/xStrato/video-encoder-golang/infrastructure/database/contexts"
)

type JobRepository struct {
	dbContext *contexts.EncoderContext
}

func NewJobRepository(ctx *contexts.EncoderContext) *JobRepository {
	return &JobRepository{ctx}
}

func (r *JobRepository) Insert(e interfaces.Entity) (interfaces.Entity, error) {
	model, ok := e.(*entities.Job)
	if !ok {
		return nil, fmt.Errorf("cannot cast '%v' as Job entity", e)
	}
	db := r.dbContext.GetDBConnection()
	if err := db.Create(&model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (r *JobRepository) Find(id string) (interfaces.Entity, error) {
	var job entities.Job
	db := r.dbContext.GetDBConnection()
	db.Preload("Video").First(&job, "id=?", id)

	if job.IsValid() != nil {
		return nil, fmt.Errorf("job ID '%v' doest not exist", id)
	}
	return &job, nil
}

func (r *JobRepository) Update(e interfaces.Entity) (interfaces.Entity, error) {
	db := r.dbContext.GetDBConnection()
	if err := db.Save(e).Error; err != nil {
		return nil, err
	}
	return e, nil
}
