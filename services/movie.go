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
//	@Summary		Adds a new movie
//	@Description	Adds a new movie with the given details including title, description, release date, and rating
//	@Tags			movie
//	@Accept			json
//	@Produce		json
//	@Param			movie	body		models.Movie	true	"Movie to add"
//	@Success		200		{object}	models.Movie	"Successfully added the movie"
//	@Failure		400		{string}	string			"Invalid request body"
//	@Failure		500		{string}	string			"Error creating movie"
//	@Router			/movie-add [post]
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

	if err := PG.db.Create(&data).Error; err != nil {
		log.Error().Err(err).Msg("Error creating actor")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	return &data, nil
}

// MovieEdit godoc
//
//	@Summary		Edits an existing movie
//	@Description	Edits a movie with the specified ID based on the given update fields such as title, description, release date, rating, and associated actors
//	@Tags			movie
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int						true	"Movie ID"
//	@Param			updates	body		map[string]interface{}	true	"Fields to update"
//	@Success		200		{object}	models.Movie			"Successfully updated the movie"
//	@Failure		400		{string}	string					"Invalid request body or movie ID"
//	@Failure		404		{string}	string					"Movie not found"
//	@Failure		500		{string}	string					"Error saving movie"
//	@Router			/movie-edit/{id} [put]
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
	if err := PG.db.Preload("Actors").First(&data, "id = ?", movieID).Error; err != nil {
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
	if err := PG.db.First(&data, "id = ?", movieID).Error; err != nil {
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
			if err := PG.db.Model(&data).Association("Actors").Append(actorsToAdd); err != nil {
				log.Error().Err(err).Msg("Failed to add actors")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil, err
			}
		}
		if len(actorsToRemove) > 0 {
			log.Debug().Interface("actorsToRemove", actorsToAdd).Msg("Removing actors to movie")
			PG.db.Model(&data).Association("Actors").Delete(actorsToRemove)
		}
	}

	if err := PG.db.Save(&data).Error; err != nil {
		log.Error().Err(err).Msg("Error saving movie")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	log.Info().Int("movieID", movieID).Msg("Movie successfully updated")
	return &data, nil
}

// MovieList godoc
// @Summary Lists all movies
// @Description Retrieves a list of all movies, including their titles, descriptions, release dates, ratings, and associated actors with sorting.
// @Tags movie
// @Produce json
// @Param sort query string false "Sort by [title|rating|releasedate], prepend '-' for descending order (default: '-rating')"
// @Success 200 {array} models.Movie "Successfully retrieved all movies"
// @Failure 500 {string} string "Error retrieving movie list"
// @Router /movie-list [get]
func (PG *Postgresql) MovieList(w http.ResponseWriter, r *http.Request) (*[]models.Movie, error) {
	log.Info().Msg("MovieList called")

	var data []models.Movie
	sortParam := r.URL.Query().Get("sort")

	// Default sorting is by rating in descending order
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

		// Check if the sort parameter is valid and set the sortOrder accordingly
		if val, ok := sortFields[sortParam]; ok {
			sortOrder = val
		}
	}

	if err := PG.db.Preload("Actors").Order(sortOrder).Find(&data).Error; err != nil {
		log.Error().Err(err).Msg("Error retrieving movie list")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	log.Info().Int("movies_count", len(data)).Msg("Movies retrieved successfully")
	return &data, nil
}

// MovieDelete godoc
//
//	@Summary		Deletes a movie
//	@Description	Deletes the movie with the specified ID, including removing all associations with actors
//	@Tags			movie
//	@Produce		json
//	@Param			id	path		int		true	"Movie ID"
//	@Success		200	{string}	string	"Successfully deleted the movie"
//	@Failure		400	{string}	string	"Invalid movie ID or URL format"
//	@Failure		500	{string}	string	"Movie not found or could not be deleted"
//	@Router			/movie-delete/{id} [delete]
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

	if err := PG.db.Exec("DELETE FROM actormovies WHERE movie_id = ?", movieID).Error; err != nil {
		log.Error().Err(err).Msg("Failed to delete associated records from the join table")
		http.Error(w, "Failed to delete associated records from the join table", http.StatusInternalServerError)
		return nil, err
	}

	if err := PG.db.Where("id = ?", movieID).Delete(&models.Movie{}).Error; err != nil {
		log.Error().Err(err).Msg("Movie not found or could not be deleted")
		http.Error(w, "Movie not found or could not be deleted", http.StatusInternalServerError)
		return nil, err
	}

	log.Info().Int("movieID", movieID).Msg("Movie deleted successfully")
	return &data, nil

}
