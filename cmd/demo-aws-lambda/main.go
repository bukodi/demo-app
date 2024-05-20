package main

import (
	"context"
	"crypto/x509"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	_ "github.com/bukodi/demo-app/pkg/domain/user"
	"github.com/bukodi/demo-app/pkg/server"
	"github.com/kr/pretty"
	"log/slog"
	"os"
	"time"
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

	now := time.Now()
	_, err := x509.SystemCertPool()
	slog.Info(fmt.Sprintf("SystemCertPool took: %s", time.Since(now)))
	if err != nil {
		slog.Error("Error loading system cert pool", "err", err)
	}

	envtxt, err := os.ReadFile("env.txt")
	if err != nil {
		slog.Error("Error reading env.txt", "err", err)
	} else {
		slog.Info(fmt.Sprintf("env.txt: %s", envtxt))
	}

	srv := server.NewServer("")
	lambdaToHttp = httpadapter.NewV2(srv.RootHandler())
	lambda.Start(Handler)
	slog.Info("Lambda finished")
}
