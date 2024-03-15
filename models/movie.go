package models

import "time"

type Movie struct {
	ID          int       `gorm:"primary_key"`
	Title       string    `gorm:"type:varchar(150);not null"`
	Description string    `gorm:"type:varchar(1000)"`
	ReleaseDate time.Time `gorm:"type:date"`
	Rating      float64   `gorm:"type:decimal(2,1)"`
	Actors      []*Actor  `gorm:"many2many:actor_movies"`
}
