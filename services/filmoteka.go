package services

import (
	"encoding/json"
	"net/http"

	"vk.com/m/dto"
	"vk.com/m/models"
	"vk.com/m/utils"
)

func (PG *Postgresql) ActorAdd(w http.ResponseWriter, r *http.Request) (*models.Actor, error) {

	var actor models.Actor

	if err := json.NewDecoder(r.Body).Decode(&actor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}

	formattedDate := utils.FormatTime(actor.DateOfBirth)
	actor.DateOfBirth = formattedDate

	if err := PG.db.Save(&actor).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	return &actor, nil

}

func (PG *Postgresql) ActorEdit(w http.ResponseWriter, r *http.Request) (*dto.Actors, error) {

	var data dto.Actors

	return &data, nil
}

func (PG *Postgresql) ActorList(w http.ResponseWriter, r *http.Request) (*dto.Actors, error) {

	var data dto.Actors

	return &data, nil
}

func (PG *Postgresql) ActorDelete(w http.ResponseWriter, r *http.Request) (*dto.Actors, error) {

	var data dto.Actors

	return &data, nil
}

func (PG *Postgresql) MovieAdd(w http.ResponseWriter, r *http.Request) (*dto.Movies, error) {

	var data dto.Movies

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil, err
	}
	defer r.Body.Close()

	if err := PG.db.Save(&data).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	return &data, nil
}

func (PG *Postgresql) MovieEdit(w http.ResponseWriter, r *http.Request) (*dto.Movies, error) {

	var data dto.Movies

	return &data, nil
}

func (PG *Postgresql) MovieFind(w http.ResponseWriter, r *http.Request) (*dto.Movies, error) {

	var data dto.Movies

	return &data, nil
}

func (PG *Postgresql) MovieList(w http.ResponseWriter, r *http.Request) (*dto.Movies, error) {

	var data dto.Movies

	return &data, nil
}

func (PG *Postgresql) MovieDelete(w http.ResponseWriter, r *http.Request) (*dto.Movies, error) {

	var data dto.Movies

	return &data, nil
}
