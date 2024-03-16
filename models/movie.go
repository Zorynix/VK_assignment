package models

type Movie struct {
	ID          int    `gorm:"primary_key"`
	Title       string `gorm:"type:varchar(150);not null"`
	Description string `gorm:"type:varchar(1000)"`
	ReleaseDate string
	Rating      float64  `gorm:"type:decimal(2,1)"`
	Actors      []*Actor `gorm:"many2many:actormovies;"`
}
