package views

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// ActorAddView handles the HTTP request to add a new actor.
// It logs the function call, attempts to add a new actor by calling the ActorAdd method on the PG interface,
// handles any errors by logging them and responding with an HTTP status code (502 Bad Gateway) to indicate a gateway error,
// and if successful, responds with the newly added actor in JSON format.
func (view *View) ActorAddView() error {

	log.Info().Msg("ActorAddView called")

	data, err := view.PG.ActorAdd(view.W, view.R)
	if err != nil {
		log.Error().Err(err).Msg("Error in ActorAdd")
		view.handleError(err, http.StatusBadGateway)
		return err
	}

	view.respondWithJSON(data)
	return nil
}

// ActorEditView handles the HTTP request to edit an existing actor's details.
// Upon being called, it logs the action, calls the ActorEdit method on the PG interface with the request data,
// checks for and handles errors similarly to ActorAddView, and responds with the updated actor details in JSON format upon success.
func (view *View) ActorEditView() error {

	log.Info().Msg("ActorEditView called")

	data, err := view.PG.ActorEdit(view.W, view.R)
	if err != nil {
		log.Error().Err(err).Msg("Error in ActorEdit")
		view.handleError(err, http.StatusBadGateway)
		return err
	}

	view.respondWithJSON(data)
	return nil
}

// ActorListView processes the HTTP request to retrieve a list of all actors.
// It logs its activation, retrieves the list of actors via the PG interface's ActorList method,
// handles potential errors by reporting them and sending an HTTP 502 status code,
// and responds with the actor list in JSON format if the retrieval is successful.
func (view *View) ActorListView() error {

	log.Info().Msg("ActorListView called")

	data, err := view.PG.ActorList(view.W, view.R)
	if err != nil {
		log.Error().Err(err).Msg("Error in ActorList")
		view.handleError(err, http.StatusBadGateway)
		return err
	}

	view.respondWithJSON(data)
	return nil
}

// ActorDeleteView manages the HTTP request to delete a specific actor.
// It initiates by logging the request, then attempts to delete the actor using the ActorDelete method on the PG interface,
// handles any encountered errors by logging and responding with an HTTP 502 status code,
// and confirms successful deletion by responding with JSON data.

func (view *View) ActorDeleteView() error {

	log.Info().Msg("ActorDeleteView called")

	data, err := view.PG.ActorDelete(view.W, view.R)
	if err != nil {
		log.Error().Err(err).Msg("Error in ActorDelete")
		view.handleError(err, http.StatusBadGateway)
		return err
	}

	view.respondWithJSON(data)
	return nil
}
