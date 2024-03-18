package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/bukodi/demo-app/pkg/server"
	"log"
)

var lambdaToHttp *httpadapter.HandlerAdapterV2

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	log.Printf("Handler called with: %+v\n", req)
	resp, err := lambdaToHttp.ProxyWithContext(ctx, req)
	log.Printf("Handler returned with: %+v, %+v\n", resp, err)
	return resp, err
}

func main() {
	log.Println("Starting Lambda log")
	lambdaToHttp = httpadapter.NewV2(server.ApiV1Mux)
	lambda.Start(Handler)
	log.Println("Lambda finished")
}
