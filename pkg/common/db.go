package common

import (
	"database/sql"
	"fmt"

	// blank import for pq lib
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

// ConnectDb function establish new db connection.
// This function does not return any error,
// but throw a fatal error when establishing db connection fails.
func ConnectDb(cfg Configuration) *sql.DB {
	log.Infof("connecting to pgsql database=%s", cfg.DBHost)
	datasource := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", datasource)
	if err != nil {
		log.Fatalf("conecting to db %s failed: %v", cfg.DBHost, err)
		return nil
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("ping to db %s failed: %v", cfg.DBHost, err)
		return nil
	}

	log.Printf("connected to db %s (%s:%s) ...", cfg.DBName, cfg.DBHost, cfg.DBPort)
	return db
}
