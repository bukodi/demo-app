package user

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bukodi/demo-app/pkg/data/dyndb"
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

func (store *userStoreDynDB) migrateTable() error {

	return dyndb.MigrateTable(context.TODO(), store.dynDbSvc, &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("email"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("password_hash"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("email"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:              aws.String(store.tableName),
		BillingMode:            types.BillingModePayPerRequest,
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{},
	})

}

func (store *userStoreDynDB) Create(user *User) error {

	//TODO implement me
	panic("implement me")
}

func (store *userStoreDynDB) List() ([]*User, error) {
	//TODO implement me
	panic("implement me")
}

func (store *userStoreDynDB) ByEmail(email string) (*User, error) {
	//TODO implement me
	panic("implement me")
}

func (store *userStoreDynDB) Delete(email string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

var _ UserStore = (*userStoreDynDB)(nil)
