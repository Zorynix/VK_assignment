package dto

import "vk.com/m/models"

type ActorsMovies struct {
	Data map[string]models.ActorMovie `gorm:"serializer:json" json:"data"`
}
