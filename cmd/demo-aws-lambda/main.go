package main

import (
	"context"
	"crypto/x509"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
	_ "github.com/bukodi/demo-app/pkg/domain/user"
	"github.com/bukodi/demo-app/pkg/server"
	"github.com/kr/pretty"

	_ "github.com/bukodi/demo-app/pkg/init_by_tags"
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

	if env, err := loadEnvFromS3(); err != nil {
		slog.Error("Error loading env from S3", "err", err)
	} else {
		slog.Info(fmt.Sprintf("env from S3: %s", string(env)))
	}

	if env, err := loadEnvFromDynDb(context.TODO()); err != nil {
		slog.Error("Error loading env from DynDb", "err", err)
	} else {
		slog.Info(fmt.Sprintf("env from DynDb: %v", env))
	}

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

func loadEnvFromS3() ([]byte, error) {
	now := time.Now()
	defer func() {
		slog.Info(fmt.Sprintf("loadEnvFromS3 took: %s", time.Since(now)))
	}()
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(opts *config.LoadOptions) error {
		opts.Region = "eu-central-1"
		return nil
	})
	if err != nil {
		return nil, err
	}
	svc := s3.NewFromConfig(cfg)

	obj, err := svc.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String("demo-app-20240702"),
		Key:    aws.String("conf/env.txt"),
	})
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func loadEnvFromDynDb(ctx context.Context) (map[string]string, error) {
	now := time.Now()
	defer func() {
		slog.Info(fmt.Sprintf("loadEnvFromDynDb took: %s", time.Since(now)))
	}()
	cfg, err := config.LoadDefaultConfig(ctx, func(opts *config.LoadOptions) error {
		opts.Region = "eu-central-1"
		return nil
	})
	if err != nil {
		return nil, err
	}
	dynDbSvc := dynamodb.NewFromConfig(cfg)

	// Perform the Scan operation
	scanResponse, err := dynDbSvc.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String("demoapp-config"),
	})
	if err != nil {
		return nil, err
	}

	// Initialize an empty slice to hold the users
	type cfgEntry struct {
		Key   string
		Value string
	}

	var cfgValues = make(map[string]string)

	// Iterate over the items in the scan response
	for _, item := range scanResponse.Items {
		e := &cfgEntry{}

		// Unmarshal the item into the User
		err := attributevalue.UnmarshalMap(item, e)
		if err != nil {
			return nil, err
		}

		cfgValues[e.Key] = e.Value
	}

	return cfgValues, nil
}
