package services

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"vk.com/m/models"
	"vk.com/m/utils"
)

type Postgresql struct {
	db *gorm.DB
}

// NewPostgreSQL creates and returns a new Postgresql instance
// This function initializes a PostgreSQL database connection using the DSN environment variable
// It sets the search path to 'vk' and automatically migrates the database schemas for Actor and Movie models
// Returns a pointer to a Postgresql struct or an error if the connection or migration fails
func NewPostgreSQL(ctx context.Context) (*Postgresql, error) {

	utils.LoadEnv()

	DSN := os.Getenv("DSN")

	conn, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})

	if err != nil {
		log.Fatal().Interface("unable to create postgresql connection pool: %v", err).Msg("")
	}

	conn = conn.Debug()

	conn.Exec("SET search_path TO vk")

	err = conn.AutoMigrate(&models.Actor{}, &models.Movie{})
	if err != nil {
		log.Fatal().Interface("unable to automigrate: %v", err).Msg("")
	}

	return &Postgresql{db: conn}, nil
}

// Ping checks the connection to the PostgreSQL database
// It verifies that the database is accessible and responding to queries
// Returns an error if the database is unreachable or not responding
func (pg *Postgresql) Ping(ctx context.Context) error {
	db, err := pg.db.DB()
	if err != nil {
		log.Fatal().Interface("unable to create postgresql connection pool: %v", err).Msg("")
	}

	return db.Ping()
}

// Close terminates the PostgreSQL database connection
// It safely closes the connection pool, freeing up resources
// Logs a fatal error if closing the connection pool fails
func (pg *Postgresql) Close() {

	db, err := pg.db.DB()
	if err != nil {
		log.Fatal().Interface("unable to create postgresql connection pool: %v", err).Msg("")
	}
	db.Close()
}
