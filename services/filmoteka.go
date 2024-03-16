package services

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"vk.com/m/models"
	"vk.com/m/utils"
)

func (PG *Postgresql) ActorAdd(w http.ResponseWriter, r *http.Request) (*models.Actor, error) {

	var data models.Actor

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	formattedDate := utils.FormatTime(data.DateOfBirth)
	data.DateOfBirth = formattedDate

	if err := PG.db.Create(&data).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	return &data, nil

}

func (PG *Postgresql) ActorEdit(w http.ResponseWriter, r *http.Request) (*models.Actor, error) {

	var data models.Actor

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return nil, nil
	}
	actorIDStr := pathParts[len(pathParts)-1]
	actorID, err := strconv.Atoi(actorIDStr)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return nil, err
	}

	if err := PG.db.Preload("Movies").First(&data, "id = ?", actorID).Error; err != nil {
		http.Error(w, "Actor not found", http.StatusNotFound)
		return nil, err
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	if err := PG.db.First(&data, "id = ?", actorID).Error; err != nil {
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
			PG.db.Model(&data).Association("Movies").Append(moviesToAdd)
		}
		if len(moviesToRemove) > 0 {
			PG.db.Model(&data).Association("Movies").Delete(moviesToRemove)
		}
	}

	if err := PG.db.Save(&data).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	return &data, nil
}

func (PG *Postgresql) ActorList(w http.ResponseWriter, r *http.Request) (*[]models.Actor, error) {
	log.Println("ActorList called")
	var data []models.Actor

	if err := PG.db.Preload("Movies").Find(&data).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	return &data, nil
}

func (PG *Postgresql) ActorDelete(w http.ResponseWriter, r *http.Request) (*models.Actor, error) {

	var data models.Actor

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return nil, errors.New("invalid URL format")
	}
	actorIDStr := pathParts[len(pathParts)-1]
	actorID, err := strconv.Atoi(actorIDStr)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return nil, err
	}

	if err := PG.db.Where("id = ?", actorID).Delete(&models.Actor{}).Error; err != nil {
		http.Error(w, "Actor not found or could not be deleted", http.StatusInternalServerError)
		return nil, err
	}

	return &data, nil
}

func (PG *Postgresql) MovieAdd(w http.ResponseWriter, r *http.Request) (*models.Movie, error) {

	var data models.Movie

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	formattedDate := utils.FormatTime(data.ReleaseDate)
	data.ReleaseDate = formattedDate

	if err := PG.db.Create(&data).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	return &data, nil
}

func (PG *Postgresql) MovieEdit(w http.ResponseWriter, r *http.Request) (*models.Movie, error) {
	var data models.Movie

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return nil, nil
	}
	movieIDStr := pathParts[len(pathParts)-1]
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return nil, err
	}

	if err := PG.db.Preload("Actors").First(&data, "id = ?", movieID).Error; err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return nil, err
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	if err := PG.db.First(&data, "id = ?", movieID).Error; err != nil {
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

	if err := PG.db.Save(&data).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
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
			PG.db.Model(&data).Association("Movies").Append(actorsToAdd)
		}
		if len(actorsToRemove) > 0 {
			PG.db.Model(&data).Association("Movies").Delete(actorsToRemove)
		}
	}

	return &data, nil
}

func (PG *Postgresql) MovieFind(w http.ResponseWriter, r *http.Request) (*models.Movie, error) {

	var data models.Movie

	return &data, nil
}

func (PG *Postgresql) MovieList(w http.ResponseWriter, r *http.Request) (*[]models.Movie, error) {
	log.Println("MovieList called")
	var data []models.Movie

	if err := PG.db.Preload("Actors").Find(&data).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	return &data, nil
}

func (PG *Postgresql) MovieDelete(w http.ResponseWriter, r *http.Request) (*models.Movie, error) {

	var data models.Movie

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return nil, errors.New("invalid URL format")
	}
	movieIDStr := pathParts[len(pathParts)-1]
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return nil, err
	}

	if err := PG.db.Where("id = ?", movieID).Delete(&models.Movie{}).Error; err != nil {
		http.Error(w, "Movie not found or could not be deleted", http.StatusInternalServerError)
		return nil, err
	}

	return &data, nil
}
