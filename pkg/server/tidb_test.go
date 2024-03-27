package server

import (
	"crypto/tls"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"os"
	"testing"
)

func TestTiDB(t *testing.T) {
	password, _ := os.LookupEnv("TIDB_PASSWORD")

	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: "gateway01.eu-central-1.prod.aws.tidbcloud.com",
	})

	db, err := sql.Open("mysql", "iFyxguC8JirCTim.root:"+password+"@tcp(gateway01.eu-central-1.prod.aws.tidbcloud.com:4000)/fortune500?tls=tidb")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT `rank`, `company_name`, `country` FROM fortune500.fortune500_2018_2022 LIMIT 100")
	if err != nil {
		t.Fatal(err)
	}
	rank, company_name, country := 0, "", ""
	for rows.Next() {
		err = rows.Scan(&rank, &company_name, &country)
		if err == nil {
			t.Logf("%d, %s, %s", rank, company_name, country)
		}
	}
}
