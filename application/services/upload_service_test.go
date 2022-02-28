package services_test

// import (
// 	"log"
// 	"os"
// 	"testing"

// 	"github.com/joho/godotenv"
// 	uuid "github.com/satori/go.uuid"
// 	"github.com/stretchr/testify/require"
// 	"github.com/xStrato/video-encoder-golang/application/services"
// 	"github.com/xStrato/video-encoder-golang/domain/entities"
// 	"github.com/xStrato/video-encoder-golang/infrastructure/database/contexts"
// 	"github.com/xStrato/video-encoder-golang/infrastructure/database/repositories"
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// func init() {
// 	err := godotenv.Load("../../.env")
// 	if err != nil {
// 		log.Fatalln("Error loading .env file")
// 	}
// }

// func TestUploadServiceMethods(t *testing.T) {

// 	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	encoderContext := contexts.NewEncoderContext(db)
// 	encoderContext.RunMigrations()
// 	videoRepo := repositories.NewVideoRepository(encoderContext)
// 	video := entities.NewVideo(uuid.NewV4().String(), "video_download_test.mp4")
// 	videoService := services.NewVideoServiceWithVideo(video, *videoRepo)

// 	videoService.Download("video-encoder-golang")
// 	videoService.Fragment()
// 	videoService.Encode()
// 	defer videoService.Finish()

// 	t.Run("ProcessUpload_ValidUploadService_ShouldNotReturnError", func(t *testing.T) {
// 		//Arrange
// 		bkt := "video-encoder-golang"
// 		path := os.Getenv("LOCAL_STORAGE_PATH") + "/" + video.ID
// 		uploadService := services.NewUploadService(bkt, path)
// 		done := make(chan string)
// 		//Act
// 		uploadService.ProcessUpload(20, done)
// 		result := <-done
// 		videoService.Finish()
// 		//Assert
// 		require.Equal(t, result, "upload completed")
// 	})
// }
