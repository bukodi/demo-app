package user

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bukodi/demo-app/pkg/data/dyndb"
	"log/slog"
	"testing"
)

// Based on this description: https://dynobase.dev/dynamodb-golang-query-examples/
// and this: https://docs.aws.amazon.com/code-library/latest/ug/go_2_dynamodb_code_examples.html
func init() {
	if !testing.Testing() {
		initDynamodbStore()
	}
}

func initDynamodbStore() error {
	cli, tablePrefix := dyndb.Client()
	if cli == nil {
		return nil
	}
	us := userStoreDynDB{
		dynDbSvc:  cli,
		tableName: tablePrefix + "users",
	}
	err := us.migrateTable()
	if err != nil {
		return err
	}
	SetUserStore(&us)
	return nil
}

type userStoreDynDB struct {
	dynDbSvc  *dynamodb.Client
	tableName string
}

func userKey(u *User) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"Email": &types.AttributeValueMemberS{Value: u.Email},
	}
}

func (store *userStoreDynDB) migrateTable() error {

	return dyndb.MigrateTable(context.TODO(), store.dynDbSvc, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("Email"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("Email"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String(store.tableName),
		BillingMode: types.BillingModePayPerRequest,
	})
}

func (store *userStoreDynDB) Create(ctx context.Context, u *User) error {
	attrs, err := attributevalue.MarshalMap(u)
	if err != nil {
		return err
	}

	resp, err := store.dynDbSvc.PutItem(ctx, &dynamodb.PutItemInput{
		Item:                   attrs,
		TableName:              aws.String(store.tableName),
		ConditionExpression:    aws.String("attribute_not_exists(Email)"),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	})
	if err != nil {
		return err
	}
	_ = resp
	return nil
}

func (store *userStoreDynDB) List(ctx context.Context) ([]*User, error) {
	//TODO implement me
	panic("implement me")
}

func (store *userStoreDynDB) ByEmail(ctx context.Context, email string) (*User, error) {
	response, err := store.dynDbSvc.GetItem(
		context.TODO(),
		&dynamodb.GetItemInput{
			Key:                    userKey(&User{Email: email}),
			TableName:              aws.String(store.tableName),
			ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		})
	if err != nil {
		return nil, err
	}

	u := User{}
	err = attributevalue.UnmarshalMap(response.Item, &u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (store *userStoreDynDB) Delete(ctx context.Context, email string) (bool, error) {
	resp, err := store.dynDbSvc.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(store.tableName),
		Key:       userKey(&User{Email: email}),
	})
	if err != nil {
		return true, err
	}

	_ = resp
	slog.Info(fmt.Sprintf("User email=%s deleted from the table", email))
	return true, err
}

var _ UserStore = (*userStoreDynDB)(nil)
