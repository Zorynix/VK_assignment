package services

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"vk.com/m/models"
	"vk.com/m/utils"
)

type Database interface {
	Ping(ctx context.Context) error

	ActorAdd(w http.ResponseWriter, r *http.Request) (*models.Actor, error)
	ActorEdit(w http.ResponseWriter, r *http.Request) (*models.Actor, error)
	ActorList(w http.ResponseWriter, r *http.Request) (*[]models.Actor, error)
	ActorDelete(w http.ResponseWriter, r *http.Request) (*models.Actor, error)

	MovieAdd(w http.ResponseWriter, r *http.Request) (*models.Movie, error)
	MovieEdit(w http.ResponseWriter, r *http.Request) (*models.Movie, error)
	MovieList(w http.ResponseWriter, r *http.Request) (*[]models.Movie, error)
	MovieFind(w http.ResponseWriter, r *http.Request) (*[]models.Movie, error)
	MovieDelete(w http.ResponseWriter, r *http.Request) (*models.Movie, error)
}

type Postgresql struct {
	DB *gorm.DB
}

// NewPostgreSQL creates and returns a new Postgresql instance
// This function initializes a PostgreSQL database connection using the DSN environment variable
// It sets the search path to 'vk' and automatically migrates the database schemas for Actor and Movie models
// Returns a pointer to a Postgresql struct or an error if the connection or migration fails
func NewPostgreSQL(ctx context.Context) (Database, error) {
	utils.LoadEnv()

	DSN := os.Getenv("DSN")
	if DSN == "" {
		return nil, errors.New("DSN is not set")
	}

	conn, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	conn = conn.Debug()

	conn.Exec("SET search_path TO vk")

	err = conn.AutoMigrate(&models.Actor{}, &models.Movie{})
	if err != nil {
		return nil, err
	}

	return &Postgresql{DB: conn}, nil
}

// Ping checks the connection to the PostgreSQL database
// It verifies that the database is accessible and responding to queries
// Returns an error if the database is unreachable or not responding
func (pg *Postgresql) Ping(ctx context.Context) error {
	DB, err := pg.DB.DB()
	if err != nil {
		log.Fatal().Interface("unable to create postgresql connection pool: %v", err).Msg("")
		return err
	}
	return DB.PingContext(ctx)
}

// Close terminates the PostgreSQL database connection
// It safely closes the connection pool, freeing up resources
// Logs a fatal error if closing the connection pool fails
func (pg *Postgresql) Close() {

	DB, err := pg.DB.DB()
	if err != nil {
		log.Fatal().Interface("unable to create postgresql connection pool: %v", err).Msg("")
	}
	DB.Close()
}
