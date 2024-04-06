package user

import (
	"context"
	"fmt"
	"github.com/bukodi/demo-app/pkg/data/gormdb"
	"gorm.io/gorm"
	"testing"
)

func init() {
	if !testing.Testing() {
		initGORMStore()
	}
}

func initGORMStore() error {
	db := gormdb.Db()
	if db == nil {
		return nil
	}
	us := userStoreGORM{
		gromDB: db,
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		return err
	}
	SetUserStore(&us)
	return nil
}

type userStoreGORM struct {
	gromDB *gorm.DB
}

var _ UserStore = (*userStoreGORM)(nil)

func (us *userStoreGORM) Create(ctx context.Context, u *User) error {
	tx := us.gromDB.Begin()
	resp := tx.Create(u)
	if resp.Error != nil {
		tx.Rollback()
		return resp.Error
	} else if resp.RowsAffected != 1 {
		tx.Rollback()
		return fmt.Errorf("record not created")
	} else {
		tx.Commit()
		return nil
	}
}

func (us *userStoreGORM) List(ctx context.Context) ([]*User, error) {
	var users []*User
	err := us.gromDB.Find(&users).Error
	return users, err
}

func (us *userStoreGORM) ByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := us.gromDB.Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (us *userStoreGORM) Delete(ctx context.Context, email string) (bool, error) {
	tx := us.gromDB.Begin()
	resp := tx.Where("email = ?", email).Delete(&User{})
	if resp.Error != nil {
		tx.Rollback()
		return false, tx.Error
	}
	tx.Commit()
	return resp.RowsAffected > 0, nil
}
