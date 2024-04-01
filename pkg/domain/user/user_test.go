package user

import (
	"os"
	"testing"
)

func TestUserStoreGORM(t *testing.T) {
	//os.Setenv("TIDB_PASSWORD", "setPassword")
	initGORMStore()

	// create a new user
	user, err := Create("email1", "password1")
	if err != nil {
		t.Fatal(err)
	}
	user2, err := VerifyPassword("email1", "password1")
	if err != nil {
		t.Fatal(err)
	}
	if user.Email != user2.Email {
		t.Fatal("email mismatch")
	}
}

func TestUserStoreDynamodb(t *testing.T) {
	os.Setenv("DYNAMODB_TABLE_PREFIX", "demoapp-")
	initDynamodbStore()

	// create a new user
	user, err := Create("email1", "password1")
	if err != nil {
		t.Fatal(err)
	}
	user2, err := VerifyPassword("email1", "password1")
	if err != nil {
		t.Fatal(err)
	}
	if user.Email != user2.Email {
		t.Fatal("email mismatch")
	}
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
