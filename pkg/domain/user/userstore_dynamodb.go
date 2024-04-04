package user

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
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

func initDynamodbStore() {
	cli := dyndb.Client()
	if cli == nil {
		return
	}
	us := userStoreDynDB{
		client: cli,
	}
	SetUserStore(&us)
}

type userStoreDynDB struct {
	client *dynamodb.Client
}

func (u *userStoreDynDB) Create(user *User) error {
	//TODO implement me
	panic("implement me")
}

func (u *userStoreDynDB) List() ([]*User, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userStoreDynDB) ByEmail(email string) (*User, error) {
	//TODO implement me
	panic("implement me")
}

var _ UserStore = (*userStoreDynDB)(nil)
