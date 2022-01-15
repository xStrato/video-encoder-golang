package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/xStrato/video-encoder-golang/domain/entities"
	"github.com/xStrato/video-encoder-golang/infrastructure/database/repositories"
)

type VideoService struct {
	Video           *entities.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (vs *VideoService) Download(bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	reader, err := client.Bucket(bucketName).Object(vs.Video.FilePath).NewReader(ctx)
	if err != nil {
		return err
	}
	defer reader.Close()
	body, err := ioutil.ReadAll(reader)

	if err != nil {
		log.Fatal(err)
	}
	path := fmt.Sprintf(os.Getenv("LOCAL_STORAGE_PATH"), "/", vs.Video.ID, ".mp4")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(body); err != nil {
		return err
	}
	log.Printf("video '%v' has been stored", vs.Video.ID)
	return nil
}
