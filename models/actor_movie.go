package models

type ActorMovie struct {
	ActorID int `gorm:"primary_key"`
	MovieID int `gorm:"primary_key"`
}
