package server

import (
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

var db *sql.DB
var db_init_err error

func init() {

	db, db_init_err = initDb()
	if db_init_err != nil {
		slog.Error("Error initializing database", "err", db_init_err)
	}

	ApiV1Mux.HandleFunc("GET /top10", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		err := listTop10Company(w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func initDb() (*sql.DB, error) {
	password, ok := os.LookupEnv("TIDB_PASSWORD")
	if !ok {
		return nil, errors.New("TIDB_PASSWORD not set")
	}

	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: "gateway01.eu-central-1.prod.aws.tidbcloud.com",
	})

	db, err := sql.Open("mysql", "iFyxguC8JirCTim.root:"+password+"@tcp(gateway01.eu-central-1.prod.aws.tidbcloud.com:4000)/fortune500?tls=tidb")
	if err != nil {
		errTxt := err.Error()
		if strings.Contains(errTxt, password) {
			errTxt = strings.ReplaceAll(errTxt, password, "********")
		}
		return nil, errors.New(errTxt)
	}

	return db, nil
}

func listTop10Company(writer io.Writer) error {
	if db_init_err != nil {
		return db_init_err
	}

	rows, err := db.Query("SELECT `rank`, `company_name`, `country` FROM fortune500.fortune500_2018_2022 LIMIT 10")
	if err != nil {
		return err
	}
	rank, company_name, country := 0, "", ""
	for rows.Next() {
		err = rows.Scan(&rank, &company_name, &country)
		if err == nil {
			fmt.Fprintf(writer, "%d, %s, %s\n", rank, company_name, country)
		}
	}

	return nil
}
