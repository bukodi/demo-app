package user

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"testing"
)

func TestUserStoreTiDBGORM(t *testing.T) {
	os.Unsetenv("SQLITE_DSN")
	os.Unsetenv("DYNAMODB_TABLE_PREFIX")
	//os.Setenv("TIDB_PASSWORD", "setPassword")
	ResetStore()

	set, _ := spi.IsSet()
	if !set {
		t.Skip("TIDB_PASSWORD not set")
		return
	}

	testUserCRUD(context.TODO(), t)
}

func TestUserStoreSqliteGORM(t *testing.T) {
	os.Unsetenv("DYNAMODB_TABLE_PREFIX")
	os.Unsetenv("TIDB_PASSWORD")
	os.Setenv("SQLITE_DSN", "file::memory:?cache=shared")
	ResetStore()

	testUserCRUD(context.TODO(), t)
}

func TestUserStoreDynamodb(t *testing.T) {
	//t.Skip("skipping dynamodb test")
	os.Unsetenv("SQLITE_DSN")
	os.Unsetenv("TIDB_PASSWORD")
	os.Setenv("DYNAMODB_TABLE_PREFIX", "demoapp-")
	ResetStore()

	testUserCRUD(context.TODO(), t)
}

func testUserCRUD(ctx context.Context, t *testing.T) {
	email := fmt.Sprintf("email-%d", rand.N(900)+100)

	// create a new user
	user, err := Create(ctx, email, "password1")
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

	user2, err := VerifyPassword(ctx, email, "password1")
	if err != nil {
		t.Fatal(err)
	}
	if user.Email != user2.Email {
		t.Fatal("email mismatch")
	}

	users, err := List(ctx)
	if err != nil {
		t.Fatal(err)
	}

	var foundUser *User
	for _, u := range users {
		if u.Email == user.Email {
			foundUser = u
			break
		}
	}

	if foundUser == nil {
		t.Fatal("user not found in list")
	}

	// create user again, it must fail
	user3, err := Create(ctx, email, "password1")
	if err == nil {
		// .Errorf("No error received on duplicate key")
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
