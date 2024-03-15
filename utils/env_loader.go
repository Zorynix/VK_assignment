package utils

import (
	"github.com/rs/zerolog/log"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Panic().Msg("---failed to load .env file---")
	}
}
