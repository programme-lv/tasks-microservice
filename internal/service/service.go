package service

import (
	"github.com/programme-lv/tasks-microservice/internal/domain"
)

type TaskRepo interface {
	GetTask(id string) (*domain.Task, error)
	ListTasks() ([]domain.Task, error)
}

type TaskService struct {
	repo TaskRepo
}

func NewTaskService(repo TaskRepo) *TaskService {
	return &TaskService{repo: repo}
}
