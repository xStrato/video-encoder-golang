package repositories

import (
	"fmt"

	"github.com/xStrato/video-encoder-golang/domain/entities"
	"github.com/xStrato/video-encoder-golang/domain/interfaces"
	"github.com/xStrato/video-encoder-golang/infrastructure/database/contexts"
)

type VideoRepository struct {
	dbContext *contexts.EncoderContext
}

func NewVideoRepository(ctx *contexts.EncoderContext) *VideoRepository {
	return &VideoRepository{ctx}
}

func (r *VideoRepository) Insert(e interfaces.Entity) (*interfaces.Entity, error) {
	model, ok := e.(*entities.Video)
	if !ok {
		return nil, fmt.Errorf("cannot cast '%v' as Category entity", e)
	}
	db := r.dbContext.GetDBConnection()
	if err := db.Create(&model).Error; err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *VideoRepository) Find(id string) (*interfaces.Entity, error) {
	return nil, nil
}
