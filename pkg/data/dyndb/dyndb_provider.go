package dyndb

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log/slog"
	"os"
	"sync"
)

var (
	dbInstance *dynamodb.Client
	once       sync.Once
)

func Client() *dynamodb.Client {
	once.Do(func() {
		db, err := initGorm()
		if err != nil {
			slog.Error("can't initialize DynamoDB: %s", err.Error(), "err", err)
		} else if db == nil {
			slog.Debug("DynamoDB isn't configured.")
		} else {
			dbInstance = db
		}
	})

	return dbInstance
}

func initGorm() (*dynamodb.Client, error) {
	_, ok := os.LookupEnv("DYNAMODB_TABLE_PREFIX")
	if !ok {
		slog.Debug("DYNAMODB_TABLE_PREFIX not set")
		return nil, nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), func(opts *config.LoadOptions) error {
		opts.Region = "eu-central-1"
		return nil
	})
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg)
	return client, nil
}
