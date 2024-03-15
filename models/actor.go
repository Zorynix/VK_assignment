package models

type Actor struct {
	ID          int    `gorm:"primary_key"`
	Name        string `gorm:"type:varchar(255);not null"`
	Gender      string `gorm:"type:char(1)"`
	DateOfBirth string
	Movies      []*Movie `gorm:"many2many:actor_movies;foreignKey:ID;joinForeignKey:ActorID;References:ID;joinReferences:MovieID"`
}
