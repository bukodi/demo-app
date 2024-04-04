package user

import (
	"github.com/bukodi/demo-app/pkg/data/gormdb"
	"gorm.io/gorm"
	"testing"
)

func init() {
	if !testing.Testing() {
		initGORMStore()
	}
}

func initGORMStore() {
	db := gormdb.Db()
	if db == nil {
		return
	}
	us := userStoreGORM{
		gromDB: db,
	}
	db.AutoMigrate(&User{})
	SetUserStore(&us)
}

type userStoreGORM struct {
	gromDB *gorm.DB
}

var _ UserStore = (*userStoreGORM)(nil)

func (us *userStoreGORM) Create(u *User) error {
	return us.gromDB.Create(u).Error
}

func (us *userStoreGORM) List() ([]*User, error) {
	var users []*User
	err := us.gromDB.Find(&users).Error
	return users, err
}

func (us *userStoreGORM) ByEmail(email string) (*User, error) {
	var u User
	err := us.gromDB.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (us *userStoreGORM) Delete(email string) (bool, error) {
	tx := us.gromDB.Where("email = ?", email).Delete(&User{})
	return tx.RowsAffected > 0, tx.Error
}
