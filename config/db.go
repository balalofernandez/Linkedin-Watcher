package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DatabaseConfiguration struct {
	Driver   string
	Dbname   string
	Username string
	Password string
	Host     string
	Port     string
	LogMode  bool
}

func DbConfiguration() string {
	dbName := viper.GetString("MASTER_DB_NAME")
	dbUser := viper.GetString("MASTER_DB_USER")
	dbPassword := viper.GetString("MASTER_DB_PASSWORD")
	dbHost := viper.GetString("MASTER_DB_HOST")
	dbPort := viper.GetString("MASTER_DB_PORT")
	dbSslMode := viper.GetString("MASTER_SSL_MODE")

	// Use proper PostgreSQL connection string format
	dbDSN := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPassword, dbHost, dbPort, dbName, dbSslMode,
	)
	return dbDSN
}
