package infected

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

// Contract is the interface will define the service methods
type Contract interface {
	StoreInfected(Request) error
}

// Client will be our base struct
type Client struct {
	db *sql.DB
}

// NewClient will return an infected client pointer
func NewClient(db *sql.DB) *Client {
	return &Client{db: db}
}

// Request will formalize the endpoint request
type Request struct {
	UUIDs   []string  `json:"uuids"`
	SinceTS time.Time `json:"since_ts"`
}

// StoreInfected will reach the github infected API and return a list of infecteds
func (ec *Client) StoreInfected(req Request) error {
	log.WithFields(log.Fields{
		"package": "infected",
		"request": req,
	}).Infoln("Requesting infected from GitHub")

	// validation checks
	if len(req.UUIDs) == 0 {
		return fmt.Errorf("No `uuids``")
	}
	if req.SinceTS.IsZero() {
		return fmt.Errorf("No `since_ts`")
	}

	// Save the request on the database
	err := ec.persistRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (ec *Client) persistRequest(req Request) error {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	// Initializing transaction
	tx, err := ec.db.Begin()
	if err != nil {
		return err
	}

	// Attempting the persist
	log.Println(string(reqBytes))

	// TODO REVIEW
	sqlStatement := "INSERT INTO had_contact_with_infected (uuid, since_ts) VALUES ($1, $2)"

	for _, uuid := range req.UUIDs {
		_, err = tx.Exec(sqlStatement, uuid, req.SinceTS)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Committing the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
