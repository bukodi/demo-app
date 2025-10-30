package gormdb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func defaultGormConfig() *gorm.Config {
	return &gorm.Config{}
}

func InitGorm(ctx context.Context) (*gorm.DB, error) {
	if db, err := initTiDBGorm(ctx); err != nil || db != nil {
		return db, err
	}
	if db, err := initSqliteGorm(ctx); err != nil || db != nil {
		return db, err
	}
	return nil, nil
}

func initSqliteGorm(ctx context.Context) (*gorm.DB, error) {
	dsn, ok := os.LookupEnv("SQLITE_DSN")
	if !ok {
		slog.Info("SQLITE_DSN not set")
		return nil, nil
	}

	gormCfg := defaultGormConfig()
	sqlDb, err := gorm.Open(sqlite.Open(dsn), gormCfg)
	if err != nil {
		return nil, err
	}
	return sqlDb, nil
}

func initTiDBGorm(ctx context.Context) (*gorm.DB, error) {
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

	gormCfg := defaultGormConfig()
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
