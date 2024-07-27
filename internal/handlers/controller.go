package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/programme-lv/tasks-microservice/internal/service"
)

type Controller struct {
	taskSrv *service.TaskService

	publicBucketCloudFrontHost string
}

func NewController(taskSrv *service.TaskService) *Controller {
	return &Controller{
		taskSrv:                    taskSrv,
		publicBucketCloudFrontHost: "dvhk4hiwp1rmf.cloudfront.net",
	}
}

func (c *Controller) RegisterRoutes(r chi.Router) {
	r.Use(middleware.Logger)

	r.Route("/tasks", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/", c.ListTasks)
			r.Get("/{id}", c.GetTask)
		})
	})

}
