package s3repo

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pelletier/go-toml/v2"
	"github.com/programme-lv/tasks-microservice/internal/domain"
)

type taskS3Repo struct {
	s3Client   *s3.Client
	taskBucket string
}

func NewTaskS3Repo(s3Client *s3.Client, taskBucket string) (*taskS3Repo, error) {
	return &taskS3Repo{
		s3Client:   s3Client,
		taskBucket: taskBucket,
	}, nil
}

func (repo *taskS3Repo) GetTask(id string) (*domain.Task, error) {
	manifest, err := repo.getTaskManifestFromS3(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task manifest: %v", err)
	}

	task, err := constructTaskFromManifest(id, manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to construct task: %v", err)
	}

	return task, nil
}

func constructTaskFromManifest(id string, manifest *TaskTomlManifest) (
	*domain.Task, error) {
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

func (repo *taskS3Repo) getTaskManifestFromS3(id string) (*TaskTomlManifest, error) {
	output, err := repo.s3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &repo.taskBucket,
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
	manifests, err := repo.listTaskManifestsFromS3()
	if err != nil {
		return nil, fmt.Errorf("failed to list task manifests: %v", err)
	}

	tasks := []domain.Task{}
	for _, manifest := range manifests {
		task, err := constructTaskFromManifest(manifest.TaskFullName, &manifest)
		if err != nil {
			return nil, fmt.Errorf("failed to construct task: %v", err)
		}

		tasks = append(tasks, *task)
	}

	return tasks, nil
}

func (repo *taskS3Repo) listTaskManifestsFromS3() ([]TaskTomlManifest, error) {
	output, err := repo.s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: &repo.taskBucket,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}

	manifests := []TaskTomlManifest{}
	for _, item := range output.Contents {
		manifest, err := repo.getTaskManifestFromS3(*item.Key)
		if err != nil {
			log.Printf("failed to get task manifest: %v", err)
			continue
		}
		manifests = append(manifests, *manifest)
	}

	return manifests, nil

}
