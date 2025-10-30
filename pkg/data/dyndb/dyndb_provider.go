package dyndb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func InitDynamoDB(ctx context.Context) (*dynamodb.Client, string, error) {
	tablePrefix, ok := os.LookupEnv("DYNAMODB_TABLE_PREFIX")
	if !ok {
		slog.Debug("DYNAMODB_TABLE_PREFIX not set")
		return nil, "", nil
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), func(opts *config.LoadOptions) error {
		opts.Region = "eu-central-1"
		return nil
	})
	if err != nil {
		return nil, tablePrefix, err
	}

	client := dynamodb.NewFromConfig(cfg)
	return client, tablePrefix, nil
}

func MigrateTable(ctx context.Context, svc *dynamodb.Client, tableDesc *dynamodb.CreateTableInput) error {
	currantDesc, err := svc.DescribeTable(ctx, &dynamodb.DescribeTableInput{
		TableName: tableDesc.TableName,
	})

	if err != nil {
		var resourceNotFound *types.ResourceNotFoundException
		if !errors.As(err, &resourceNotFound) {
			return fmt.Errorf("Failed to describe table %s, %w.\n", *tableDesc.TableName, err)
		} else {
			// Create the table
			resp, err := svc.CreateTable(ctx, tableDesc)
			if err != nil {
				return fmt.Errorf("Failed to create table %s, %w.\n", *tableDesc.TableName, err)
			}
			slog.Info(fmt.Sprintf("DynamoDB table %s created with arn: %s", *resp.TableDescription.TableName, *resp.TableDescription.TableArn))
			return nil
		}
	}

	// Check if the table stuctuire meets the requirements
	if len(currantDesc.Table.KeySchema) != len(tableDesc.KeySchema) {
		return errors.New("table key schema does not match")
	}
	for i, expectedKey := range tableDesc.KeySchema {
		currentKey := currantDesc.Table.KeySchema[i]
		if *currentKey.AttributeName != *expectedKey.AttributeName || currentKey.KeyType != expectedKey.KeyType {
			return fmt.Errorf("table key schema does not match (expected:%v actual:%v)", expectedKey, currentKey)
		}
	}
	if len(currantDesc.Table.AttributeDefinitions) != len(tableDesc.AttributeDefinitions) {
		return errors.New("table attribute definitions does not match")
	}
	for i, expectedAttr := range tableDesc.AttributeDefinitions {
		currentAttr := currantDesc.Table.AttributeDefinitions[i]
		if *currentAttr.AttributeName != *expectedAttr.AttributeName || currentAttr.AttributeType != expectedAttr.AttributeType {
			return fmt.Errorf("table attribute definition does not match (expected:%v actual:%v)", expectedAttr, currentAttr)
		}
	}
	return nil
}
