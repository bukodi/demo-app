package data

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type UserStore interface {
	// Create inserts a new user into the store
	Create(user User) error
	// List returns all users from the store
	List() ([]User, error)
	// Verify checks if the given email and password is valid
	Verify(email, password string) error
}

type userStoreGORM struct {
	gromDB *gorm.DB
}

func (us *userStoreGORM) Create(email string, password string) error {
	user := new(User)
	user.Email = email
	pswHash, err := hashAndSalt([]byte(password))
	if err != nil {
		return err
	}
	user.PasswordHash = pswHash
	return us.gromDB.Create(user).Error
}

func (us *userStoreGORM) List() ([]User, error) {
	var users []User
	err := us.gromDB.Find(&users).Error
	return users, err
}

var ErrInvalidCredentials = fmt.Errorf("invalid credentials")

func (us *userStoreGORM) Verify(email, password string) error {
	var user User
	err := us.gromDB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return err
	}

	if comparePasswords(user.PasswordHash, []byte(password)) {
		return ErrInvalidCredentials
	}
	return nil
}

func hashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func NewUserStoreGORM(db *gorm.DB) UserStore {
	return nil // &userStoreGORM{gromDB: db}
}
