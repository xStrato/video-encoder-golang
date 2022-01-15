package entities_test

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/xStrato/video-encoder-golang/domain/entities"
)

func TestVideoStructValidations(t *testing.T) {

	t.Run("IsValid_InvalidDefaultStruct_ShouldReturnError", func(t *testing.T) {
		//Arrange
		video := entities.Video{}
		//Act
		err := video.IsValid()
		//Assert
		require.Error(t, err)
	})

	t.Run("IsValid_InvalidUUID_ShouldReturnError", func(t *testing.T) {
		//Arrange
		video := entities.Video{}
		video.ID = "abc"
		//Act
		err := video.IsValid()
		//Assert
		require.EqualError(t, err, "id: abc does not validate as uuid")
	})

	t.Run("IsValid_ValidStructConstruction_ShouldNotReturnError", func(t *testing.T) {
		//Arrange
		video := entities.NewVideo(uuid.NewV4().String(), "<file_path>")
		//Act
		err := video.IsValid()
		//Assert
		require.Nil(t, err)
	})
}
