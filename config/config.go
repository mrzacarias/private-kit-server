package config

import (
	"github.com/spf13/viper"
)

// AppConfig will contain private-kit-server configurable information
type AppConfig struct {
	Port        string
	MetricsPort string

	// Database
	DBHost     string
	DBPort     string
	DBDatabase string
	DBUser     string
	DBPassword string
	DBSSLMode  string
	DBDriver   string
}

func init() {
	viper.SetEnvPrefix("PKS")
	viper.AutomaticEnv()
	viper.SetDefault("port", "8080")
	viper.SetDefault("metrics_port", "8081")

	// Database
	viper.SetDefault("db_host", "test_db")
	viper.SetDefault("db_port", "5432")
	viper.SetDefault("db_database", "private-kit-server_test")
	viper.SetDefault("db_user", "postgres")
	viper.SetDefault("db_password", "postgres")
	viper.SetDefault("db_ssl_mode", "disable")
	viper.SetDefault("db_driver", "postgres")
}

// GetConfig will generate the standard AppConfig
func GetConfig() AppConfig {
	return AppConfig{
		Port:        viper.GetString("port"),
		MetricsPort: viper.GetString("metrics_port"),

		// Database
		DBHost:     viper.GetString("db_host"),
		DBPort:     viper.GetString("db_port"),
		DBDatabase: viper.GetString("db_database"),
		DBUser:     viper.GetString("db_user"),
		DBPassword: viper.GetString("db_password"),
		DBSSLMode:  viper.GetString("db_ssl_mode"),
		DBDriver:   viper.GetString("db_driver"),
	}
}
