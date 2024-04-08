package routes

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"

	"vk.com/m/services"
)

type Router struct {
	PG services.Database
}

func Routes(addr *string) {
	postgres, err := services.NewPostgreSQL(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize PostgreSQL")
	}

	router := Router{PG: postgres}

	router.V1Routes()
	log.Info().Msgf("Starting server on port %d...", 8000)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal().Err(err).Msg("Cannot start HTTP server")
	}
}
