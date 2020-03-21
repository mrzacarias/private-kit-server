package testing

import (
	"database/sql"
	"fmt"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres" // need for CloudSQL proxy driver
	"github.com/mrzacarias/private-kit-server/config"
	_ "github.com/lib/pq" // import the postgres driver
	log "github.com/sirupsen/logrus"
)

// NewTestDBClient will initialize the Test Database and return the client
func NewTestDBClient() *sql.DB {
	cfg := config.GetConfig()

	log.WithFields(log.Fields{"Host": cfg.DBHost}).Infoln("Initializing Test Database!")

	// Attempting connection
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=5432 user=postgres password=postgres dbname=private-kit-server_test sslmode=disable", cfg.DBHost))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Error while setting up testing DB")
	}

	// Pinging to ensure that we are on
	if err := db.Ping(); err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Error pinging testing database connection")
	}

	return db
}

// Truncate - removes all rows from test DB
func Truncate(db *sql.DB) (sql.Result, error) {
	return db.Exec("TRUNCATE had_contact_with_infected")
}
