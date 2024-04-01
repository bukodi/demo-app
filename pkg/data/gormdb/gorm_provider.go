package gormdb

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log/slog"
	"os"
	"strings"
	"sync"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

func Db() *gorm.DB {
	once.Do(func() {
		db, err := initGorm()
		if err != nil {
			slog.Error("can't initialize GORM db: %s", err.Error(), "err", err)
		} else if db == nil {
			slog.Debug("GORM db isn't configured.")
		} else {
			dbInstance = db
		}
	})

	return dbInstance
}

func initGorm() (*gorm.DB, error) {
	password, ok := os.LookupEnv("TIDB_PASSWORD")
	if !ok {
		slog.Debug("TIDB_PASSWORD not set")
		return nil, nil
	}

	tidb_user := "iFyxguC8JirCTim.root"
	tidb_host := "gateway01.eu-central-1.prod.aws.tidbcloud.com"
	tidb_port := "4000"
	tidb_db_name := "fortune500"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=true", tidb_user, password, tidb_host, tidb_port, tidb_db_name)

	gormCfg := &gorm.Config{}
	sqlDb, err := gorm.Open(mysql.Open(dsn), gormCfg)
	if err != nil {
		errTxt := err.Error()
		if strings.Contains(errTxt, password) {
			// Prevent error containing password
			errTxt = strings.ReplaceAll(errTxt, password, "********")
			err = errors.New(errTxt)
		}
		return nil, err
	}
	return sqlDb, nil
}
