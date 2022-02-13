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

func NewVideoServiceWithVideo(v *entities.Video, repo repositories.VideoRepository) VideoService {
	return VideoService{v, repo}
}
func NewVideoService(repo *repositories.VideoRepository) *VideoService {
	return &VideoService{
		videoRepository: *repo,
	}
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
	source := fmt.Sprint(os.Getenv("LOCAL_STORAGE_PATH"), "/", vs.video.ID, ".mp4")
	target := fmt.Sprint(os.Getenv("LOCAL_STORAGE_PATH"), "/", vs.video.ID, ".frag")
	cmd := exec.CommandContext(context.Background(), "mp4fragment", source, target)
	return execCommand(cmd)
}

func (vs *VideoService) Encode() error {
	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, fmt.Sprint(os.Getenv("LOCAL_STORAGE_PATH"), "/", vs.video.ID, ".frag"))
	cmdArgs = append(cmdArgs, "--use-segment-timeline")
	cmdArgs = append(cmdArgs, "-o")
	cmdArgs = append(cmdArgs, fmt.Sprint(os.Getenv("LOCAL_STORAGE_PATH"), "/", vs.video.ID))
	cmdArgs = append(cmdArgs, "-f")
	cmdArgs = append(cmdArgs, "--exec-dir")
	cmdArgs = append(cmdArgs, "/opt/bento4/bin/")

	cmd := exec.CommandContext(context.Background(), "mp4dash", cmdArgs...)
	return execCommand(cmd)
}
func (vs *VideoService) Finish() error {
	err := os.Remove(fmt.Sprint(os.Getenv("LOCAL_STORAGE_PATH"), "/", vs.video.ID, ".mp4"))
	if err != nil {
		log.Println("error removing mp4: ", vs.video.ID, ".mp4")
		return err
	}
	err = os.Remove(fmt.Sprint(os.Getenv("LOCAL_STORAGE_PATH"), "/", vs.video.ID, ".frag"))
	if err != nil {
		log.Println("error removing frag: ", vs.video.ID, ".frag")
		return err
	}
	err = os.RemoveAll(fmt.Sprint(os.Getenv("LOCAL_STORAGE_PATH"), "/", vs.video.ID))
	if err != nil {
		log.Println("error removing dir: ", vs.video.ID)
		return err
	}
	log.Println("files have been removed at:", vs.video.ID)
	return nil
}
func (vs *VideoService) InsertVideo() error {
	if _, err := vs.videoRepository.Insert(vs.video); err != nil {
		return err
	}
	return nil
}
func execCommand(cmd *exec.Cmd) error {
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	if len(output) > 0 {
		log.Printf("=====> Output: %s\n", string(output))
	}
	return nil
}
