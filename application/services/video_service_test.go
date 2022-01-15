package services_test

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/xStrato/video-encoder-golang/application/services"
	"github.com/xStrato/video-encoder-golang/domain/entities"
	"github.com/xStrato/video-encoder-golang/infrastructure/database/contexts"
	"github.com/xStrato/video-encoder-golang/infrastructure/database/repositories"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
}

func TestVideoServiceMethods(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	encoderContext := contexts.NewEncoderContext(db)
	encoderContext.RunMigrations()
	// jobRepo := repositories.NewJobRepository(encoderContext)
	videoRepo := repositories.NewVideoRepository(encoderContext)
	video := entities.NewVideo(uuid.NewV4().String(), "game-progress-mechanics.mp4")

	t.Run("Download_ValidJVideoServiceStruct_ShouldNotReturnError", func(t *testing.T) {
		//Arrange
		videoService := services.NewVideoService(video, *videoRepo)
		//Act
		err := videoService.Download("video-encoder-golang")
		//Assert
		require.Nil(t, err)
	})
}
