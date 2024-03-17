package models

// Actor represents an actor in the movie database.
// It contains information about the actor's ID, name, gender, date of birth, and the movies they've acted in.
// The struct uses GORM annotations to define how it maps to your database schema, specifying field properties like primary keys and field types.
//
// Fields:
// - ID: The unique identifier for the actor. It's marked as the primary key in the database.
// - Name: The name of the actor, stored as a varchar(255) in the database and cannot be null.
// - Gender: The gender of the actor, stored as a single character (M or F) indicating male or female, respectively.
// - DateOfBirth: The date of birth of the actor, stored as a string. No specific type is enforced in the database schema through GORM annotations.
// - Movies: A slice of pointers to Movie structs, representing the many-to-many relationship between actors and movies. This is managed through the "actormovies" join table.
type Actor struct {
	ID          int    `gorm:"primary_key"`
	Name        string `gorm:"type:varchar(255);not null"`
	Gender      string `gorm:"type:char(1)"`
	DateOfBirth string
	Movies      []*Movie `gorm:"many2many:actormovies;"`
}
