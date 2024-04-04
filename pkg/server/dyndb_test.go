package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func TestCreateTable(t *testing.T) {
	t.Skip("skipping test")
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(opts *config.LoadOptions) error {
		opts.Region = "eu-central-1"
		return nil
	})
	if err != nil {
		panic(err)
	}

	svc := dynamodb.NewFromConfig(cfg)

	tableName := "my-table4"

	desc, err := svc.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: aws.String(tableName),
	})

	t.Logf("Table description: %v", desc)
	if err != nil {
		var resourceNotFound *types.ResourceNotFoundException
		if errors.As(err, &resourceNotFound) {
			fmt.Printf("Table %s does not exist.\n", tableName)
		} else {
			t.Fatalf("Failed to describe table %s, %v.\n", tableName, err)
		}
	} else {
		t.Logf("Table %s exists.\n", tableName)
	}

	out, err := svc.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:              aws.String(tableName),
		BillingMode:            types.BillingModePayPerRequest,
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(out)
}

func MigrateTable(ctx *context.Context, svc *dynamodb.Client, tableDesc *dynamodb.CreateTableInput) error {
	currantDesc, err := svc.DescribeTable(context.TODO(), &dynamodb.DescribeTableInput{
		TableName: tableDesc.TableName,
	})

	if err != nil {
		var resourceNotFound *types.ResourceNotFoundException
		if !errors.As(err, &resourceNotFound) {
			return fmt.Errorf("Failed to describe table %s, %w.\n", *tableDesc.TableName, err)
		} else {
			// Create the table
			resp, err := svc.CreateTable(*ctx, tableDesc)
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
