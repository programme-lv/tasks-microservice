package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/programme-lv/tasks-microservice/internal/handlers"
	"github.com/programme-lv/tasks-microservice/internal/repositories/ddbtaskrepo"
	"github.com/programme-lv/tasks-microservice/internal/service"
)

const taskTable = "ProglvTasks"

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-central-1"))
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}
	dynamodbClient := dynamodb.NewFromConfig(cfg)

	repo := ddbtaskrepo.NewDynamoDbTaskRepo(dynamodbClient, taskTable)

	taskService := service.NewTaskService(repo)
	controller := handlers.NewController(taskService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	controller.RegisterRoutes(r)

	fmt.Println("Server started at port 8080")
	fmt.Println(http.ListenAndServe(":8080", r))
}
