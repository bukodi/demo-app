package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/bukodi/demo-app/pkg/server"
	"github.com/kr/pretty"
	"log"
	"log/slog"
	"os"
)

var lambdaToHttp *httpadapter.HandlerAdapterV2

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Printf("Handler called with: %s\n", pretty.Sprintf("%+v", req))
	resp, err := lambdaToHttp.ProxyWithContext(ctx, req)
	log.Printf("Handler returned with: %+v, %+v\n", pretty.Sprintf("%+v", resp), err)
	return resp, err
}

func main() {
	slog.Info("Starting Lambda slog")
	log.Println("Starting Lambda log")

	envtxt, err := os.ReadFile("env.txt")
	if err != nil {
		log.Printf("Error reading env.txt: %s", err.Error())
	} else {
		log.Printf("env.txt: %s", envtxt)
	}

	lambdaToHttp = httpadapter.NewV2(server.ApiV1Mux)
	lambda.Start(Handler)
	log.Println("Lambda finished")
	slog.Info("Lambda finished slog")
}
