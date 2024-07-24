package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/programme-lv/tasks-microservice/internal/handlers"
	"github.com/programme-lv/tasks-microservice/internal/service"

	awschi "github.com/awslabs/aws-lambda-go-api-proxy/chi"
)

func main() {
	taskService := service.NewTaskService(nil)
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
