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

func NewPostgreSQL(ctx context.Context) (*Postgresql, error) {

	utils.LoadEnv()

	DSN := os.Getenv("DSN")

	conn, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})

	if err != nil {
		log.Fatal().Interface("unable to create postgresql connection pool: %v", err).Msg("")
	}

	err = conn.AutoMigrate(&models.Actor{}, &models.Movie{}, &models.ActorMovie{})
	if err != nil {
		log.Fatal().Interface("unable to automigrate: %v", err).Msg("")
	}

	return &Postgresql{db: conn}, nil
}

func (pg *Postgresql) Ping(ctx context.Context) error {
	db, err := pg.db.DB()
	if err != nil {
		log.Fatal().Interface("unable to create postgresql connection pool: %v", err).Msg("")
	}

	return db.Ping()
}

func (pg *Postgresql) Close() {

	db, err := pg.db.DB()
	if err != nil {
		log.Fatal().Interface("unable to create postgresql connection pool: %v", err).Msg("")
	}
	db.Close()
}
