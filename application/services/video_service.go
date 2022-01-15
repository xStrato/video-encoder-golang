package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
	"github.com/xStrato/video-encoder-golang/domain/entities"
	"github.com/xStrato/video-encoder-golang/infrastructure/database/repositories"
)

type VideoService struct {
	video           *entities.Video
	videoRepository repositories.VideoRepository
}

func NewVideoService(v *entities.Video, repo repositories.VideoRepository) VideoService {
	return VideoService{v, repo}
}

func (vs *VideoService) Download(bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	reader, err := client.Bucket(bucketName).Object(vs.video.FilePath).NewReader(ctx)
	if err != nil {
		return err
	}
	defer reader.Close()
	body, err := ioutil.ReadAll(reader)

	if err != nil {
		log.Fatal(err)
	}
	path := os.Getenv("LOCAL_STORAGE_PATH") + "/" + vs.video.ID + ".mp4"
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.Write(body); err != nil {
		return err
	}
	log.Printf("video '%v' has been stored", vs.video.ID)
	return nil
}

func (vs *VideoService) Fragment() error {
	path := os.Getenv("LOCAL_STORAGE_PATH") + "/" + vs.video.ID
	if err := os.Mkdir(path, os.ModePerm); err != nil {
		return err
	}

	source := fmt.Sprintf(os.Getenv("LOCAL_STORAGE_PATH"), "/", vs.video.ID, ".mp4")
	target := fmt.Sprintf(os.Getenv("LOCAL_STORAGE_PATH"), "/", vs.video.ID, ".frag")
	cmd := exec.CommandContext(context.Background(), "mp4fragment", source, target)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	if len(output) > 0 {
		log.Printf("=====> Output: %s\n", string(output))
	}
	return nil
}
