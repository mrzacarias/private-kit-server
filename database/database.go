package database

import (
	"database/sql"
	"fmt"

	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres" // need for CloudSQL proxy driver
	"github.com/mrzacarias/private-kit-server/config"
	_ "github.com/lib/pq" // import the postgres driver
	log "github.com/sirupsen/logrus"
)

// NewDBClient will initialize the Database and return the client
func NewDBClient() *sql.DB {
	cfg := config.GetConfig()

	log.WithFields(log.Fields{"Host": cfg.DBHost, "Port": cfg.DBPort, "Database": cfg.DBDatabase}).Infoln("Initializing Database!")

	// Attempting connection
	db, err := sql.Open(cfg.DBDriver, fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBDatabase,
		cfg.DBSSLMode,
	))
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Error while setting up DB")
	}

	// Pinging to ensure that we are on
	if err := db.Ping(); err != nil {
		log.WithFields(log.Fields{"error": err}).Fatal("Error pinging database connection")
	}

	return db
}
