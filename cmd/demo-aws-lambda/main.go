package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	"github.com/bukodi/demo-app/pkg/server"
)

var lambdaToHttp *httpadapter.HandlerAdapterV2

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return lambdaToHttp.ProxyWithContext(ctx, req)
}

func main() {
	lambdaToHttp = httpadapter.NewV2(server.ApiV1Mux)
	lambda.Start(Handler)
}
