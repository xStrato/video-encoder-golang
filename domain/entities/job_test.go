package entities_test

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/xStrato/video-encoder-golang/domain/entities"
)

func TestJobStructValidations(t *testing.T) {

	t.Run("IsValid_InvalidDefaultStruct_ShouldReturnError", func(t *testing.T) {
		//Arrange
		job := entities.Job{}
		//Act
		err := job.IsValid()
		//Assert
		require.Error(t, err)
	})

	t.Run("IsValid_InvalidUUID_ShouldReturnError", func(t *testing.T) {
		//Arrange
		job := entities.Job{}
		job.ID = "abc"
		//Act
		err := job.IsValid()
		//Assert
		require.EqualError(t, err, "id: abc does not validate as uuid")
	})

	t.Run("IsValid_ValidStructConstruction_ShouldNotReturnError", func(t *testing.T) {
		//Arrange
		video := entities.NewVideo(uuid.NewV4().String(), "<file_path>")
		job := entities.NewJob("<output_path>", "Converted", video)
		//Act
		err1 := video.IsValid()
		err2 := job.IsValid()
		//Assert
		require.Nil(t, err1)
		require.Nil(t, err2)
	})
}
