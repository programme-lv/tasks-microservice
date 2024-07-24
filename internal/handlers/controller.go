package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/programme-lv/tasks-microservice/internal/service"
)

type Controller struct {
	TaskSrv *service.TaskService
}

func NewController(taskSrv *service.TaskService) *Controller {
	return &Controller{TaskSrv: taskSrv}
}

func (c *Controller) RegisterRoutes(r chi.Router) {
	r.Use(middleware.Logger)

	r.Route("/tasks", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/", c.ListTasks)
			r.Get("/{id}", c.GetTask)
		})
	})

	r.Get("/task-pdfs/{sha}", c.GetPdf)

}
