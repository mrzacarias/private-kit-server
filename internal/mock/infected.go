package mock

import (
	"fmt"
	"os"

	"github.com/mrzacarias/private-kit-server/internal/infected"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

// --------------- EMOJI MOCK ---------------

// InfectedClient will mock when an infected was found
type InfectedClient struct{}

// StoreInfected will mock InfectedClient StoreInfected
func (ec *InfectedClient) StoreInfected(req infected.Request) error {
	if req.UUIDs[0] == "ef872b78-3df5-483c-a764-6f33e1a23898" {
		log.WithFields(log.Fields{"mock": true, "package": "infected"}).Info("StoreInfected worked")
		return nil
	}
	log.WithFields(log.Fields{"mock": true, "package": "infected"}).Info("StoreInfected didn't worked")
	return fmt.Errorf("FAIL")
}
