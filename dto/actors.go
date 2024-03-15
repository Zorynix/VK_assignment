package dto

import "vk.com/m/models"

type Actors struct {
	Data map[string]models.Actor `gorm:"serializer:json" json:"data"`
}
