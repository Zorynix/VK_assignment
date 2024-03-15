package dto

import "vk.com/m/models"

type Movies struct {
	Data map[string]models.Movie `gorm:"serializer:json" json:"data"`
}
