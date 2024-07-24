package service

import (
	"github.com/programme-lv/tasks-microservice/internal/domain"
)

func (x *TaskService) GetTask(id string) (*domain.Task, error) {
	return x.repo.GetTask(id)
}

func (x *TaskService) ListTasks() ([]domain.Task, error) {
	return x.repo.ListTasks()
}
