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
	videoRepo := repositories.NewVideoRepository(encoderContext)
	video := entities.NewVideo(uuid.NewV4().String(), "video_download_test.mp4")
	videoService := services.NewVideoService(video, *videoRepo)

	t.Run("Download_ValidJVideoServiceStruct_ShouldNotReturnError", func(t *testing.T) {
		//Arrange
		//Act
		err := videoService.Download("video-encoder-golang")
		//Assert
		require.Nil(t, err)
	})

	t.Run("Fragment_ValidJVideoServiceStruct_ShouldNotReturnError", func(t *testing.T) {
		//Arrange
		//Act
		err := videoService.Fragment()
		//Assert
		require.Nil(t, err)
	})
}
