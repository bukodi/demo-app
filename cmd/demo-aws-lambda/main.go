package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/bukodi/demo-app/pkg/server"
	"github.com/kr/pretty"
	"log/slog"
	"net/http"
	"os"
)

var lambdaToHttp *httpadapter.HandlerAdapterV2

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	slog.Debug(fmt.Sprintf("Handler called with: %s\n", pretty.Sprint(req)))
	resp, err := lambdaToHttp.ProxyWithContext(ctx, req)
	slog.Debug(fmt.Sprintf("Handler returned with: %s, %+v\n", pretty.Sprint(resp), err))
	return resp, err
}

func main() {
	slog.Info("Starting Lambda")

	envtxt, err := os.ReadFile("env.txt")
	if err != nil {
		slog.Error("Error reading env.txt", "err", err)
	} else {
		slog.Info(fmt.Sprintf("env.txt: %s", envtxt))
	}

	rootMux := http.NewServeMux()
	rootMux.Handle("/api/v1/", http.StripPrefix("/api/v1", server.ApiV1Mux))

	lambdaToHttp = httpadapter.NewV2(rootMux)
	lambda.Start(Handler)
	slog.Info("Lambda finished")
}
