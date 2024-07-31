package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/programme-lv/tasks-microservice/internal/handlers"
	"github.com/programme-lv/tasks-microservice/internal/repositories/ddbtaskrepo"
	"github.com/programme-lv/tasks-microservice/internal/service"

	awschi "github.com/awslabs/aws-lambda-go-api-proxy/chi"
)

func main() {
	taskService := service.NewTaskService(getDynamoDbRepo())
	controller := handlers.NewController(taskService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(corsMiddleware)

	controller.RegisterRoutes(r)

	chiLambda := awschi.NewV2(r)

	handler := func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (
		events.APIGatewayV2HTTPResponse, error) {
		return chiLambda.ProxyWithContextV2(ctx, req)
	}

	lambda.Start(handler)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods",
			"GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getDynamoDbRepo() service.TaskRepo {
	tableName := os.Getenv("TASKS_TABLE_NAME")
	if tableName == "" {
		panic("TASKS_TABLE_NAME environment variable is not set")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("eu-central-1"))
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}
	dynamoClient := dynamodb.NewFromConfig(cfg)
	repo := ddbtaskrepo.NewDynamoDbTaskRepo(dynamoClient,
		tableName)
	return repo
}
