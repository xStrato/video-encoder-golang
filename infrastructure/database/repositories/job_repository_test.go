package repositories_test

import (
	"fmt"
	"log"
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/xStrato/video-encoder-golang/domain/entities"
	"github.com/xStrato/video-encoder-golang/infrastructure/database/contexts"
	"github.com/xStrato/video-encoder-golang/infrastructure/database/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestJobRepositoryMethods(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	encoderContext := contexts.NewEncoderContext(db)
	encoderContext.RunMigrations()
	jobRepo := repositories.NewJobRepository(encoderContext)
	videoRepo := repositories.NewVideoRepository(encoderContext)

	var insertedJobUUID string

	t.Run("Insert_ValidJobStruct_ShouldNotReturnError", func(t *testing.T) {
		//Arrange
		video := entities.NewVideo(uuid.NewV4().String(), "<file_path>")
		videoRepo.Insert(video)

		job := entities.NewJob("<output>", "Pending", video)
		insertedJobUUID = job.GetId()
		//Act
		entity, err := jobRepo.Insert(job)
		j, _ := entity.(*entities.Job)
		//Assert
		require.NotNil(t, entity)
		require.Nil(t, err)
		require.Equal(t, job.ID, entity.GetId())
		require.Equal(t, video.GetId(), j.VideoID)
		require.IsType(t, job, entity)
	})

	t.Run("Find_ValidInsertedJobUUID_ShouldNotReturnError", func(t *testing.T) {
		//Arrange
		//Act
		entity, err := jobRepo.Find(insertedJobUUID)
		//Assert
		require.NotNil(t, entity)
		require.Nil(t, err)
		require.Equal(t, insertedJobUUID, entity.GetId())
	})

	t.Run("Update_ValidJobStruct_ShouldNotReturnError", func(t *testing.T) {
		//Arrange
		video := entities.NewVideo(uuid.NewV4().String(), "<file_path>")
		videoRepo.Insert(video)

		job := entities.NewJob("<output>", "Pending", video)
		entity, err := jobRepo.Insert(job)
		j, _ := entity.(*entities.Job)
		//Act
		j.Status = "Completed"
		jobRepo.Update(j)
		//Assert
		require.NotNil(t, j)
		require.Nil(t, err)
		require.Equal(t, job.ID, j.GetId())
		require.Equal(t, job.Status, j.Status)
		require.IsType(t, job, j)
	})

	t.Run("Find_InvalidInsertedJobUUID_ShouldReturnError", func(t *testing.T) {
		//Arrange
		invalidID := uuid.NewV4().String()
		errMsg := fmt.Sprintf("job ID '%v' doest not exist", invalidID)
		//Act
		entity, err := jobRepo.Find(invalidID)
		//Assert
		require.Nil(t, entity)
		require.EqualError(t, err, errMsg)
	})
}
