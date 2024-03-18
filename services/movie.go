package services

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"vk.com/m/models"
	"vk.com/m/utils"
)

// MovieAdd godoc
//
// @Security ApiKeyAuth
// @SecurityRequirement ApiKeyAuth
// @Summary Adds a new movie
// @Description Adds a new movie with the given details including title, description, release date, and rating. Requires 'admin' role.
// @Tags movie
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer [JWT token]"
// @Param movie body models.Movie true "Movie to add"
// @Success 200 {object} models.Movie "Successfully added the movie"
// @Failure 400 "Invalid request body"
// @Failure 401 "Unauthorized or Invalid token"
// @Failure 403 "Forbidden - Role not allowed"
// @Failure 500 "Error creating movie"
// @Router /v1/movie-add [post]
func (PG *Postgresql) MovieAdd(w http.ResponseWriter, r *http.Request) (*models.Movie, error) {

	log.Info().Msg("MovieAdd called")

	var data models.Movie

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Error().Err(err).Msg("Error decoding request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	formattedDate := utils.FormatTime(data.ReleaseDate)
	data.ReleaseDate = formattedDate

	if err := PG.DB.Create(&data).Error; err != nil {
		log.Error().Err(err).Msg("Error creating actor")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	return &data, nil
}

// MovieEdit godoc
//
// @Security ApiKeyAuth
// @SecurityRequirement ApiKeyAuth
// @Summary Edits an existing movie
// @Description Edits a movie with the specified ID based on the given update fields such as title, description, release date, rating, and associated actors. Requires 'admin' role.
// @Tags movie
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer [JWT token]"
// @Param id path int true "Movie ID"
// @Param updates body map[string]interface{} true "Fields to update"
// @Success 200 {object} models.Movie "Successfully updated the movie"
// @Failure 400 "Invalid request body or movie ID"
// @Failure 401 "Unauthorized or Invalid token"
// @Failure 403 "Forbidden - Role not allowed"
// @Failure 404 "Movie not found"
// @Failure 500 "Error saving movie"
// @Router /v1/movie-edit/{id} [put]
func (PG *Postgresql) MovieEdit(w http.ResponseWriter, r *http.Request) (*models.Movie, error) {

	log.Info().Msg("MovieEdit called")

	var data models.Movie

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		log.Error().Msg("Invalid URL format")
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return nil, nil
	}
	movieIDStr := pathParts[len(pathParts)-1]
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid movie ID")
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return nil, err
	}

	log.Debug().Int("movieID", movieID).Msg("Fetching movie from database")
	if err := PG.DB.Preload("Actors").First(&data, "id = ?", movieID).Error; err != nil {
		log.Error().Err(err).Msg("Movie not found")
		http.Error(w, "Movie not found", http.StatusNotFound)
		return nil, err
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Error().Err(err).Msg("Error decoding updates")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	log.Debug().Interface("updates", updates).Msg("Applying updates to movie")
	if err := PG.DB.First(&data, "id = ?", movieID).Error; err != nil {
		log.Error().Err(err).Msg("Movie not found for updating")
		http.Error(w, "Movie not found", http.StatusNotFound)
		return nil, err
	}

	for field, value := range updates {
		switch field {
		case "title":
			if title, ok := value.(string); ok {
				data.Title = title
			}
		case "description":
			if description, ok := value.(string); ok {
				data.Description = description
			}
		case "releasedate":
			if releasedatestr, ok := value.(string); ok {
				formattedReleaseDate := utils.FormatTime(releasedatestr)
				data.ReleaseDate = formattedReleaseDate
			}
		case "rating":
			if rating, ok := value.(float64); ok {
				data.Rating = rating
			}
		}
	}

	if actorIDsInterface, ok := updates["actors"].([]interface{}); ok {
		var actorsToAdd []models.Actor
		var currentActorIDs []int

		for _, m := range data.Actors {
			currentActorIDs = append(currentActorIDs, m.ID)
		}

		for _, idInterface := range actorIDsInterface {
			idInt, err := utils.InterfaceToInt(idInterface)
			if err == nil && !utils.Contains(currentActorIDs, idInt) {
				actorsToAdd = append(actorsToAdd, models.Actor{ID: idInt})
			}
		}

		var actorsToRemove []models.Actor
		for _, currentID := range currentActorIDs {
			if !utils.ContainsInterfaceAsInt(actorIDsInterface, currentID) {
				actorsToRemove = append(actorsToRemove, models.Actor{ID: currentID})
			}
		}

		if len(actorsToAdd) > 0 {
			if err := PG.DB.Model(&data).Association("Actors").Append(actorsToAdd); err != nil {
				log.Error().Err(err).Msg("Failed to add actors")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil, err
			}
		}
		if len(actorsToRemove) > 0 {
			log.Debug().Interface("actorsToRemove", actorsToAdd).Msg("Removing actors to movie")
			PG.DB.Model(&data).Association("Actors").Delete(actorsToRemove)
		}
	}

	if err := PG.DB.Save(&data).Error; err != nil {
		log.Error().Err(err).Msg("Error saving movie")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	log.Info().Int("movieID", movieID).Msg("Movie successfully updated")
	return &data, nil
}

// MovieList godoc
//
// @Security ApiKeyAuth
// @SecurityRequirement ApiKeyAuth
// @Summary Lists all movies
// @Description Retrieves a list of all movies, including their titles, descriptions, release dates, ratings, and associated actors with sorting. Available to both 'admin' and 'user' roles.
// @Tags movie
// @Produce json
// @Param Authorization header string true "Bearer [JWT token]"
// @Param sort query string false "Sort by [title|rating|releasedate], prepend '-' for descending order (default: '-rating')"
// @Success 200 {array} models.Movie "Successfully retrieved all movies"
// @Failure 401 "Unauthorized or Invalid token"
// @Failure 403 "Forbidden - Role not allowed"
// @Failure 500 "Error retrieving movie list"
// @Router /v1/movie-list [get]
func (PG *Postgresql) MovieList(w http.ResponseWriter, r *http.Request) (*[]models.Movie, error) {
	log.Info().Msg("MovieList called")

	var data []models.Movie
	sortParam := r.URL.Query().Get("sort")

	sortOrder := "rating DESC"
	if sortParam != "" {
		// Map query parameters to database columns
		sortFields := map[string]string{
			"title":        "title",
			"-title":       "title DESC",
			"rating":       "rating",
			"-rating":      "rating DESC",
			"releasedate":  "release_date",
			"-releasedate": "release_date DESC",
		}

		if val, ok := sortFields[sortParam]; ok {
			sortOrder = val
		}
	}

	if err := PG.DB.Preload("Actors").Order(sortOrder).Find(&data).Error; err != nil {
		log.Error().Err(err).Msg("Error retrieving movie list")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	log.Info().Int("movies_count", len(data)).Msg("Movies retrieved successfully")
	return &data, nil
}

// MovieFind godoc
//
// @Security ApiKeyAuth
// @SecurityRequirement ApiKeyAuth
// @Summary Searches for movies by title or actor name
// @Description Searches for movies by a fragment of the title or by a fragment of an actor's name. Available to both 'admin' and 'user' roles.
// @Tags movie
// @Produce json
// @Param Authorization header string true "Bearer [JWT token]"
// @Param title query string false "Fragment of the movie title"
// @Param actor query string false "Fragment of the actor's name"
// @Success 200 {array} models.Movie "Successfully found movies"
// @Failure 401 "Unauthorized or Invalid token"
// @Failure 403 "Forbidden - Role not allowed"
// @Failure 500 "Error retrieving movie list"
// @Router /v1/movie-find [get]
func (PG *Postgresql) MovieFind(w http.ResponseWriter, r *http.Request) (*[]models.Movie, error) {

	log.Info().Msg("MovieFind called")

	var movies []models.Movie
	title := r.URL.Query().Get("title")
	actor := r.URL.Query().Get("actor")

	query := PG.DB.Model(&models.Movie{})

	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}

	if actor != "" {
		query = query.Joins("JOIN actormovies ON actormovies.movie_id = movies.id").
			Joins("JOIN actors ON actors.id = actormovies.actor_id").
			Where("actors.name ILIKE ?", "%"+actor+"%")
	}

	if err := query.Preload("Actors").Find(&movies).Error; err != nil {
		log.Error().Err(err).Msg("Error searching for movies")
		http.Error(w, "Error searching for movies", http.StatusInternalServerError)
		return nil, err
	}

	return &movies, nil
}

// MovieDelete godoc
//
// @Security ApiKeyAuth
// @SecurityRequirement ApiKeyAuth
// @Summary Deletes a movie
// @Description Deletes the movie with the specified ID, including removing all associations with actors. Requires 'admin' role.
// @Tags movie
// @Produce json
// @Param Authorization header string true "Bearer [JWT token]"
// @Param id path int true "Movie ID"
// @Success 200 "Successfully deleted the movie"
// @Failure 400 "Invalid movie ID or URL format"
// @Failure 401 "Unauthorized or Invalid token"
// @Failure 403 "Forbidden - Role not allowed"
// @Failure 500 "Movie not found or could not be deleted"
// @Router /v1/movie-delete/{id} [delete]
func (PG *Postgresql) MovieDelete(w http.ResponseWriter, r *http.Request) (*models.Movie, error) {

	log.Info().Msg("MovieDelete called")

	var data models.Movie

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		log.Error().Msg("Invalid URL format")
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return nil, errors.New("invalid URL format")
	}
	movieIDStr := pathParts[len(pathParts)-1]
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid movie ID")
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return nil, err
	}

	if err := PG.DB.Exec("DELETE FROM actormovies WHERE movie_id = ?", movieID).Error; err != nil {
		log.Error().Err(err).Msg("Failed to delete associated records from the join table")
		http.Error(w, "Failed to delete associated records from the join table", http.StatusInternalServerError)
		return nil, err
	}

	if err := PG.DB.Where("id = ?", movieID).Delete(&models.Movie{}).Error; err != nil {
		log.Error().Err(err).Msg("Movie not found or could not be deleted")
		http.Error(w, "Movie not found or could not be deleted", http.StatusInternalServerError)
		return nil, err
	}

	log.Info().Int("movieID", movieID).Msg("Movie deleted successfully")
	return &data, nil

}
