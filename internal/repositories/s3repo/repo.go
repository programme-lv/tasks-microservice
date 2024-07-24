package s3repo

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pelletier/go-toml/v2"
	"github.com/programme-lv/tasks-microservice/internal/domain"
)

var (
	testfilesBucket = "proglv-tests"
	taskBucket      = "proglv-tasks"
	pdfBucket       = "proglv-pdfs"
)

type taskS3Repo struct {
	s3Client *s3.Client
}

func NewTaskS3Repo() (*taskS3Repo, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("eu-central-1"))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	s3Client := s3.NewFromConfig(cfg)
	return &taskS3Repo{s3Client: s3Client}, nil
}

func (repo *taskS3Repo) GetTask(id string) (*domain.Task, error) {
	manifest, err := getTaskManifestFromS3(repo, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task manifest: %v", err)
	}

	task, err := domain.NewTask(id, manifest.TaskFullName)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %v", err)
	}

	task.SetCpuTimeLimitSecs(manifest.CpuTimeInSecs)
	task.SetDifficulty(manifest.Difficulty)
	task.SetMemoryLimitMBytes(manifest.MemoryLimMB)
	task.SetOriginOlympiad(manifest.OriginOlympiad)
	task.SetProblemTags(manifest.ProblemTags)
	task.SetTaskFullName(manifest.TaskFullName)

	pdfs := []domain.PdfSha256Ref{}
	for _, pdf := range manifest.PDFSHA256s {
		pdfs = append(pdfs, domain.PdfSha256Ref{
			Language: pdf.Language,
			Sha256:   pdf.SHA256,
		})
	}

	task.SetPdfStatement(pdfs)

	return task, nil
}

func getTaskManifestFromS3(repo *taskS3Repo, id string) (*TaskTomlManifest, error) {
	output, err := repo.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &taskBucket,
		Key:    &id,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download object: %v", err)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(output.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %v", err)
	}

	manifest := TaskTomlManifest{}
	err = toml.Unmarshal(buf.Bytes(), &manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal manifest: %v", err)
	}

	return &manifest, nil
}

func (repo *taskS3Repo) ListTasks() ([]domain.Task, error) {
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
