package models

// Movie represents a movie in the database.
// It includes details about the movie's ID, title, description, release date, and rating, as well as the actors who have appeared in the movie.
// GORM annotations are used to specify the database schema details, such as primary keys and field types.
//
// Fields:
// - ID: The unique identifier for the movie, serving as the primary key in the database.
// - Title: The title of the movie, stored as a varchar(150) in the database and marked as not nullable.
// - Description: A description of the movie, allowing for up to varchar(1000) characters. This field is not marked as not null, so it's optional.
// - ReleaseDate: The release date of the movie, stored as a string. Like the DateOfBirth in the Actor model, no specific database type is enforced via GORM.
// - Rating: The movie's rating, stored as a decimal with one digit after the decimal point (e.g., 8.5). This allows for a rating scale of 0.0 to 9.9.
// - Actors: A slice of pointers to Actor structs, indicating the many-to-many relationship with actors through the "actormovies" join table. This shows which actors have appeared in the movie.
type Movie struct {
	ID          int    `gorm:"primary_key"`
	Title       string `gorm:"type:varchar(150);not null"`
	Description string `gorm:"type:varchar(1000)"`
	ReleaseDate string
	Rating      float64  `gorm:"type:decimal(2,1)"`
	Actors      []*Actor `gorm:"many2many:actormovies;"`
}
