package services

import (
	"context"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"cloud.google.com/go/storage"
)

type UploadService struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewUploadService(output, path string) *UploadService {
	return &UploadService{
		OutputBucket: output,
		VideoPath:    path,
	}
}

func (v *UploadService) UploadObject(objPath string, client *storage.Client, ctx context.Context) error {
	path := strings.Split(objPath, os.Getenv("LOCAL_STORAGE_PATH")+"/")

	f, err := os.Open(objPath)
	if err != nil {
		return err
	}
	defer f.Close()
	wc := client.Bucket(v.OutputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err = io.Copy(wc, f); err != nil {
		return err
	}
	if err = wc.Close(); err != nil {
		return err
	}
	return nil
}

func (v *UploadService) loadPaths() error {
	err := filepath.Walk(v.VideoPath, func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			v.Paths = append(v.Paths, path)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
func (v *UploadService) getUploadClient() (*storage.Client, context.Context, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, nil, err
	}
	return client, ctx, nil
}

func (v *UploadService) ProcessUpload(concurrency int, done chan string) error {
	in := make(chan int, runtime.NumCPU())
	out := make(chan string)

	if err := v.loadPaths(); err != nil {
		return err
	}

	client, ctx, err := v.getUploadClient()
	if err != nil {
		return err
	}
	for p := 0; p < concurrency; p++ {
		go v.worker(in, out, client, ctx)
	}

	go func() {
		for i := 0; i < len(v.Paths); i++ {
			in <- i
		}
		close(in)
	}()

	for resp := range out {
		if resp != "" {
			done <- resp
			break
		}
	}
	return nil
}

func (v *UploadService) worker(in chan int, out chan string, client *storage.Client, ctx context.Context) {
	for index := range in {
		err := v.UploadObject(v.Paths[index], client, ctx)
		if err != nil {
			v.Errors = append(v.Errors, v.Paths[index])
			log.Panicf("an error occured during upload: %v. Message: %v", v.Paths[index], err)
			out <- err.Error()
		}
		out <- ""
	}
	out <- "upload completed"
}
