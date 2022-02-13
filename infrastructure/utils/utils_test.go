package utils_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xStrato/video-encoder-golang/infrastructure/utils"
)

func TestUtilsFunctions(t *testing.T) {

	t.Run("IsJson_ValidJsonString_ShouldReturnNil", func(t *testing.T) {
		//Arrange
		js := ` {
			"id": "28dcb8aa-b730-4d25-b1bb-6b906ce89789",
			"file_path": "video.mp4",
			"status": "pending"
		}`
		//Act
		err := utils.IsJson(js)
		//Assert
		require.Nil(t, err)
	})

	t.Run("IsJson_InvalidJsonString_ShouldReturnError", func(t *testing.T) {
		//Arrange
		js := ` {
			id: "28dcb8aa-b730-4d25-b1bb-6b906ce89789",
			file_path: "video.mp4",
			status: "pending"
		}`
		//Act
		err := utils.IsJson(js)
		//Assert
		require.Error(t, err)
	})
}
