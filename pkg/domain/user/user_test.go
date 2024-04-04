package user

import (
	"context"
	"os"
	"testing"
)

func TestUserStoreGORM(t *testing.T) {
	os.Unsetenv("DYNAMODB_TABLE_PREFIX")
	//os.Setenv("TIDB_PASSWORD", "setPassword")
	if err := initGORMStore(); err != nil {
		t.Fatal(err)
	} else if !IsUserStoreSet() {
		t.Skip("store not initialized")
	}

	ctx := context.TODO()
	t.Run("GORM-TIDBServerless", func(t *testing.T) {
		testUserCRUD(ctx, t)
	})
}
func TestUserStoreDynamodb(t *testing.T) {
	//t.Skip("skipping dynamodb test")
	os.Unsetenv("TIDB_PASSWORD")
	os.Setenv("DYNAMODB_TABLE_PREFIX", "demoapp-")
	if err := initDynamodbStore(); err != nil {
		t.Fatal(err)
	} else if !IsUserStoreSet() {
		t.Skip("store not initialized")
	}

	ctx := context.TODO()
	t.Run("DynamoDB", func(t *testing.T) {
		testUserCRUD(ctx, t)
	})
}

func testUserCRUD(ctx context.Context, t *testing.T) {
	// create a new user
	user, err := Create(ctx, "email1", "password1")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if deleted, err := store().Delete(ctx, user.Email); err != nil {
			t.Fatal(err)
		} else if !deleted {
			t.Fatal("user not deleted")
		}
	}()

	user2, err := VerifyPassword(ctx, "email1", "password1")
	if err != nil {
		t.Fatal(err)
	}
	if user.Email != user2.Email {
		t.Fatal("email mismatch")
	}

	// create user again, it must fail
	user3, err := Create(ctx, "email1", "password1")
	if err == nil {
		t.Errorf("No error received on duplicate key")
	} else {
		t.Logf("expected error: %v", err)
	}
	_ = user3
}

func TestPasswordHash(t *testing.T) {
	plainPassword := "Passw0rd"
	saltedPsw1, err := HashAndSaltPassword(plainPassword)
	if err != nil {
		t.Fatal(err)
	}
	if CheckPassword(saltedPsw1, plainPassword) == false {
		t.Fatal("password verification failed")
	}
	saltedPsw2, err := HashAndSaltPassword(plainPassword)
	if err != nil {
		t.Fatal(err)
	}
	if CheckPassword(saltedPsw2, plainPassword) == false {
		t.Fatal("password verification failed")
	}
	if saltedPsw1 == saltedPsw2 {
		t.Fatal("salted passwords should be different")
	}
}
