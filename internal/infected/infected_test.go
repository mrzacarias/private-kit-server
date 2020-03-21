package infected_test

import (
	"database/sql"
	"testing"
	"time"

	dbTest "github.com/mrzacarias/private-kit-server/database/testing"
	"github.com/mrzacarias/private-kit-server/internal/infected"
)

var TestDB *sql.DB
var InfectedClient infected.Contract

// Helper for setup/teardown
func setup(t *testing.T) func(t *testing.T) {
	// Preparing clients
	TestDB = dbTest.NewTestDBClient()
	InfectedClient = infected.NewClient(TestDB)

	// Returning teardown function
	return func(t *testing.T) {
		// Truncating DB
		_, err := dbTest.Truncate(TestDB)
		if err != nil {
			t.Fatal("Error truncating the Database: ", err)
		}
		TestDB.Close()
	}
}

func TestStoreInfected(test *testing.T) {
	test.Run("Infected requested works", func(t *testing.T) {
		teardown := setup(test)
		defer teardown(test)

		err := InfectedClient.StoreInfected(infected.Request{UUIDs: []string{"ef872b78-3df5-483c-a764-6f33e1a23898"}, SinceTS: time.Now()})
		if err != nil {
			t.Fatal("Error on StoreInfected: ", err)
		}
	})

	test.Run("Infected requested do not work", func(t *testing.T) {
		teardown := setup(test)
		defer teardown(test)

		err := InfectedClient.StoreInfected(infected.Request{})
		if err == nil {
			t.Fatal("Should have returned an error")
		}
	})

	test.Run("Infected requests are persisted", func(t *testing.T) {
		teardown := setup(test)
		defer teardown(test)

		err := InfectedClient.StoreInfected(infected.Request{UUIDs: []string{"ef872b78-3df5-483c-a764-6f33e1a23898"}, SinceTS: time.Now()})
		if err != nil {
			t.Fatal("Error on StoreInfected: ", err)
		}

		res, err := TestDB.Exec("SELECT * FROM had_contact_with_infected")
		if err != nil {
			t.Fatal("Error on when loading infected requests from DB: ", err)
		}
		rows, err := res.RowsAffected()
		if err != nil {
			t.Fatal("Error on RowsAffected: ", err)
		}

		if rows != 1 {
			t.Fatalf("# Rows should be 1, but it was: %d", rows)
		}
	})
}
