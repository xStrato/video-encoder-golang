package repositories_test

import (
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

func TestVideoRepositoryMethods(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	encoderContext := contexts.NewEncoderContext(db)
	encoderContext.RunMigrations()
	videoRepo := repositories.NewVideoRepository(encoderContext)

	var insertedVideoUUID string

	t.Run("Insert_ValidVideoStruct_ShouldNotReturnError", func(t *testing.T) {
		//Arrange
		video := entities.NewVideo(uuid.NewV4().String(), "<file_path>")
		insertedVideoUUID = video.GetId()
		//Act
		entity, err := videoRepo.Insert(video)
		//Assert
		require.NotNil(t, entity)
		require.Nil(t, err)
		require.Equal(t, video.ID, entity.GetId())
		require.IsType(t, video, entity)
	})

	t.Run("Find_ValidInsertedVideoUUID_ShouldNotReturnError", func(t *testing.T) {
		//Arrange
		//Act
		entity, err := videoRepo.Find(insertedVideoUUID)
		//Assert
		require.NotNil(t, entity)
		require.Nil(t, err)
		require.IsType(t, insertedVideoUUID, entity.GetId())
	})
}
