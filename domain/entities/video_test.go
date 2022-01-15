package entities_test

import (
	"testing"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/xStrato/video-encoder-golang/domain/entities"
)

func TestStructValidations(t *testing.T) {

	t.Run("IsValid_InvalidDefaultStruct_ShouldReturnError", func(t *testing.T) {
		video := entities.Video{}
		err := video.IsValid()
		require.Error(t, err)
	})

	t.Run("IsValid_InvalidUUID_ShouldReturnError", func(t *testing.T) {
		video := entities.Video{}
		video.ID = "abc"
		err := video.IsValid()
		require.Errorf(t, err, "abc")
	})

	t.Run("IsValid_ValidStructConstruction_ShouldNotReturnError", func(t *testing.T) {
		video := entities.NewVideo(uuid.NewV4().String(), "<file_path>")
		err := video.IsValid()
		require.Nil(t, err)
	})
}
