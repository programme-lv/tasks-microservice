package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/programme-lv/tasks-microservice/internal/handlers"
	"github.com/programme-lv/tasks-microservice/internal/repositories/s3repo"
	"github.com/programme-lv/tasks-microservice/internal/service"
)

const testBucket = "proglv-tests"

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-central-1"))
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}
	s3Client := s3.NewFromConfig(cfg)

	s3repo.NewTaskS3Repo(s3Client, testBucket)
	taskService := service.NewTaskService(nil)
	controller := handlers.NewController(taskService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	controller.RegisterRoutes(r)

	fmt.Println("Server started at port 8080")
	fmt.Println(http.ListenAndServe(":8080", r))
}
