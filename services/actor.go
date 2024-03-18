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

// ActorAdd godoc
//
// @Security ApiKeyAuth
// @SecurityRequirement ApiKeyAuth
// @Summary Adds a new actor
// @Description Adds a new actor with the given details. Requires 'admin' role.
// @Tags actor
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer [JWT token]"
// @Param actor body models.Actor true "Actor to add"
// @Success 200 {object} models.Actor "Successfully added the actor"
// @Failure 400 "Invalid request body"
// @Failure 401 "Unauthorized or Invalid token"
// @Failure 403 "Forbidden - Role not allowed"
// @Failure 500 "Error creating actor"
// @Router /v1/actor-add [post]
func (PG *Postgresql) ActorAdd(w http.ResponseWriter, r *http.Request) (*models.Actor, error) {

	log.Info().Msg("ActorAdd called")

	var data models.Actor

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Error().Err(err).Msg("Error decoding request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	formattedDate := utils.FormatTime(data.DateOfBirth)
	data.DateOfBirth = formattedDate

	if err := PG.DB.Create(&data).Error; err != nil {
		log.Error().Err(err).Msg("Error creating actor")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	log.Info().Msg("Actor added successfully")
	return &data, nil

}

// ActorEdit godoc
//
// @Security ApiKeyAuth
// @SecurityRequirement ApiKeyAuth
// @Summary Edits an existing actor
// @Description Edits an actor with the specified ID based on the given update fields. Requires 'admin' role.
// @Tags actor
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer [JWT token]"
// @Param id path int true "Actor ID"
// @Param updates body map[string]interface{} true "Fields to update"
// @Success 200 {object} models.Actor "Successfully updated the actor"
// @Failure 400 "Invalid request body or actor ID"
// @Failure 401 "Unauthorized or Invalid token"
// @Failure 403 "Forbidden - Role not allowed"
// @Failure 404 "Actor not found"
// @Failure 500 "Failed to save actor"
// @Router /v1/actor-edit/{id} [put]
func (PG *Postgresql) ActorEdit(w http.ResponseWriter, r *http.Request) (*models.Actor, error) {
	log.Info().Msg("ActorEdit called")

	var data models.Actor

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		log.Error().Msg("Invalid URL format")
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return nil, nil
	}
	actorIDStr := pathParts[len(pathParts)-1]
	actorID, err := strconv.Atoi(actorIDStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid actor ID")
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return nil, err
	}

	log.Debug().Int("actorID", actorID).Msg("Fetching actor from database")
	if err := PG.DB.Preload("Movies").First(&data, "id = ?", actorID).Error; err != nil {
		log.Error().Err(err).Msg("Actor not found")
		http.Error(w, "Actor not found", http.StatusNotFound)
		return nil, err
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	log.Debug().Interface("updates", updates).Msg("Applying updates to actor")
	if err := PG.DB.First(&data, "id = ?", actorID).Error; err != nil {
		log.Error().Err(err).Msg("Actor not found for updating")
		http.Error(w, "Actor not found", http.StatusNotFound)
		return nil, err
	}

	for field, value := range updates {
		switch field {
		case "name":
			if name, ok := value.(string); ok {
				data.Name = name
			}
		case "gender":
			if gender, ok := value.(string); ok && (gender == "M" || gender == "F") {
				data.Gender = gender
			}
		case "dateOfBirth":
			if dobStr, ok := value.(string); ok {
				formattedDOB := utils.FormatTime(dobStr)
				data.DateOfBirth = formattedDOB
			}
		}
	}

	if movieIDsInterface, ok := updates["movies"].([]interface{}); ok {
		var moviesToAdd []models.Movie
		var currentMovieIDs []int

		for _, m := range data.Movies {
			currentMovieIDs = append(currentMovieIDs, m.ID)
		}

		for _, idInterface := range movieIDsInterface {
			idInt, err := utils.InterfaceToInt(idInterface)
			if err == nil && !utils.Contains(currentMovieIDs, idInt) {
				moviesToAdd = append(moviesToAdd, models.Movie{ID: idInt})
			}
		}

		var moviesToRemove []models.Movie
		for _, currentID := range currentMovieIDs {
			if !utils.ContainsInterfaceAsInt(movieIDsInterface, currentID) {
				moviesToRemove = append(moviesToRemove, models.Movie{ID: currentID})
			}
		}

		if len(moviesToAdd) > 0 {
			log.Debug().Interface("moviesToAdd", moviesToAdd).Msg("Adding movies to actor")
			PG.DB.Model(&data).Association("Movies").Append(moviesToAdd)
		}
		if len(moviesToRemove) > 0 {
			log.Debug().Interface("moviesToRemove", moviesToRemove).Msg("Removing movies from actor")
			PG.DB.Model(&data).Association("Movies").Delete(moviesToRemove)
		}
	}

	if err := PG.DB.Save(&data).Error; err != nil {
		log.Error().Err(err).Msg("Failed to save actor")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	log.Info().Msg("Actor updated successfully")
	return &data, nil
}

// ActorList godoc
//
// @Security ApiKeyAuth
// @SecurityRequirement ApiKeyAuth
// @Summary Lists all actors
// @Description Retrieves a list of all actors, including their associated movies. Available to both 'admin' and 'user' roles.
// @Tags actor
// @Produce json
// @Param Authorization header string true "Bearer [JWT token]"
// @Success 200 {array} models.Actor "Successfully retrieved all actors"
// @Failure 401 "Unauthorized or Invalid token"
// @Failure 403 "Forbidden - Role not allowed"
// @Failure 500 "Error retrieving actors"
// @Router /v1/actor-list [get]
func (PG *Postgresql) ActorList(w http.ResponseWriter, r *http.Request) (*[]models.Actor, error) {
	log.Info().Msg("ActorList called")

	var data []models.Actor

	if err := PG.DB.Preload("Movies").Find(&data).Error; err != nil {
		log.Error().Err(err).Msg("Error retrieving actors")

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	log.Info().Int("count", len(data)).Msg("Successfully retrieved actors")

	return &data, nil
}

// ActorDelete godoc
//
// @Security ApiKeyAuth
// @SecurityRequirement ApiKeyAuth
// @Summary Deletes an actor
// @Description Deletes the actor with the specified ID, including removing all associated movies. Requires 'admin' role.
// @Tags actor
// @Produce json
// @Param Authorization header string true "Bearer [JWT token]"
// @Param id path int true "Actor ID"
// @Success 200 "Successfully deleted the actor"
// @Failure 400 "Invalid actor ID or URL format"
// @Failure 401 "Unauthorized or Invalid token"
// @Failure 403 "Forbidden - Role not allowed"
// @Failure 500 "Actor not found or could not be deleted"
// @Router /v1/actor-delete/{id} [delete]
func (PG *Postgresql) ActorDelete(w http.ResponseWriter, r *http.Request) (*models.Actor, error) {

	log.Info().Msg("ActorDelete called")

	var data models.Actor

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		log.Error().Msg("Invalid URL format")
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return nil, errors.New("invalid URL format")
	}
	actorIDStr := pathParts[len(pathParts)-1]
	actorID, err := strconv.Atoi(actorIDStr)
	if err != nil {
		log.Error().Err(err).Msg("Invalid actor ID")
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return nil, err
	}

	if err := PG.DB.Exec("DELETE FROM actormovies WHERE actor_id = ?", actorID).Error; err != nil {
		log.Error().Err(err).Msg("Failed to delete associated records from the join table")
		http.Error(w, "Failed to delete associated records from the join table", http.StatusInternalServerError)
		return nil, err
	}

	if err := PG.DB.Where("id = ?", actorID).Delete(&models.Actor{}).Error; err != nil {
		log.Error().Err(err).Msg("Actor not found or could not be deleted")
		http.Error(w, "Actor not found or could not be deleted", http.StatusInternalServerError)
		return nil, err
	}

	log.Info().Int("actorID", actorID).Msg("Actor successfully deleted")
	return &data, nil
}
