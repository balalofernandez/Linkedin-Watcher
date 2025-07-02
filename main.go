// @title LinkedIn Watcher API
// @version 1.0
// @description A Go-based application that monitors LinkedIn connections and automates networking actions.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"linkedin-watcher/config"
	"linkedin-watcher/db"
	_ "linkedin-watcher/docs"
	"linkedin-watcher/infra/logger"
	"linkedin-watcher/internal/routers"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func initDB(ctx context.Context, connectionString string) (db.Querier, func()) {
	// Try to initialize database, but don't fail if it's not available
	if err := db.Init(connectionString); err != nil {
		logger.Warnf("Database initialization failed: %v", err)
		logger.Infof("Starting application without database connection")
		return nil, func() {}
	}

	conn, err := pgxpool.New(ctx, connectionString)
	if err != nil {
		logger.Warnf("Database connection failed: %v", err)
		logger.Infof("Starting application without database connection")
		return nil, func() {}
	}

	queries := db.New(conn)
	logger.Infof("Database connection established successfully")
	return queries, conn.Close
}

func main() {

	//set timezone
	viper.SetDefault("SERVER_TIMEZONE", "Europe/Madrid")
	loc, _ := time.LoadLocation(viper.GetString("SERVER_TIMEZONE"))
	time.Local = loc

	if err := config.SetupConfig(); err != nil {
		logger.Fatalf("config SetupConfig() error: %s", err)
	}

	dbDSN := config.DbConfiguration()

	ctx := context.Background()
	queries, cleanup := initDB(ctx, dbDSN)
	defer cleanup()
	_ = queries // TODO: Use queries for database operations

	router := routers.SetupRoute()
	logger.Infof("Starting server on %s", config.ServerConfig())
	logger.Fatalf("%v", router.Run(config.ServerConfig()))

}
