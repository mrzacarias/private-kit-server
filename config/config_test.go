package config_test

import (
	"testing"

	cfg "github.com/mrzacarias/private-kit-server/config"
)

func TestMain(t *testing.T) {
	config := cfg.GetConfig()

	checkConfig(t, "Port", config.Port, "8080")
	checkConfig(t, "MetricsPort", config.MetricsPort, "8081")

	checkConfig(t, "DBHost", config.DBHost, "test_db")
	checkConfig(t, "DBPort", config.DBPort, "5432")
	checkConfig(t, "DBDatabase", config.DBDatabase, "private-kit-server_test")
	checkConfig(t, "DBUser", config.DBUser, "postgres")
	checkConfig(t, "DBPassword", config.DBPassword, "postgres")
	checkConfig(t, "DBSSLMode", config.DBSSLMode, "disable")
	checkConfig(t, "DBDriver", config.DBDriver, "postgres")
}

// DRY helpers for checking a config attribute
func checkConfig(t *testing.T, attr string, got string, want string) {
	if want != got {
		t.Fatalf("Attribute '%s' should be %s, but it was %s", attr, want, got)
	}
}
