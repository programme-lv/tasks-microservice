package repositories

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/klauspost/compress/zstd"
	"github.com/programme-lv/fs-task-format-parser/pkg/fstaskparser"
	"github.com/programme-lv/tasks-microservice/internal/domain"
)

const (
	testfilesBucket = "proglv-tests"
	taskBucket      = "proglv-tasks"
	pdfBucket       = "proglv-pdfs"
)

type TaskS3Repo struct {
	s3Client *s3.Client
}

func NewTaskS3Repo() (*TaskS3Repo, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	return &TaskS3Repo{s3Client: s3Client}, nil
}

func (repo *TaskS3Repo) GetTask(id string) (*domain.Task, error) {
	bucket := taskBucket
	key := id

	output, err := repo.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download object: %v", err)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(output.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %v", err)
	}

	task, err := fstaskparser.Parse(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to parse task: %v", err)
	}

	return &domain.Task{
		Id:                id,
		TaskFullName:      task.GetTaskName(),
		MemoryLimitMBytes: task.GetMemoryLimitInMegabytes(),
		CpuTimeLimitSecs:  task.GetCPUTimeLimitInSeconds(),
		Difficulty:        task.GetDifficultyOneToFive(),
		OriginOlympiad:    task.GetOriginOlympiad(),
		ProblemTags:       task.GetProblemTags(),
		PdfStatements:     []domain.PdfSha256Ref{},
	}, nil
}

func (repo *TaskS3Repo) ListTasks() ([]domain.Task, error) {
	bucket := taskBucket

	output, err := repo.s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &bucket,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}

	tasks := []domain.Task{}
	for _, item := range output.Contents {
		task, err := repo.GetTask(*item.Key)
		if err != nil {
			log.Printf("failed to get task: %v", err)
			continue
		}
		tasks = append(tasks, *task)
	}

	return tasks, nil
}
