package views

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// MovieAddView deals with the HTTP request to add a new movie.
// Similar to the actor-related views, it logs the request, adds a movie via the MovieAdd method,
// handles errors with logging and an HTTP 502 response, and returns the added movie in JSON format on success.
func (view *View) MovieAddView() error {

	log.Info().Msg("MovieAddView called")

	data, err := view.PG.MovieAdd(view.W, view.R)
	if err != nil {
		log.Error().Err(err).Msg("Error in MovieAdd")
		view.handleError(err, http.StatusBadGateway)
		return err
	}

	view.respondWithJSON(data)
	return nil
}

// MovieEditView handles the HTTP request to edit details of an existing movie.
// It logs the invocation of the function, then attempts to edit a movie's details by calling the MovieEdit method on the PG interface,
// using data from the HTTP request. If an error occurs during this process, it logs the error, responds with a 502 Bad Gateway status,
// indicating an issue with processing the request, and returns the error. If the movie is successfully edited, it responds with the updated
// movie details in JSON format.
func (view *View) MovieEditView() error {

	log.Info().Msg("MovieEditView called")

	data, err := view.PG.MovieEdit(view.W, view.R)
	if err != nil {
		log.Error().Err(err).Msg("Error in MovieEdit")
		view.handleError(err, http.StatusBadGateway)
		return err
	}

	view.respondWithJSON(data)
	return nil
}

// MovieListView manages the HTTP request to list all movies stored in the database.
// It begins by logging its execution, then retrieves the list of all movies through the MovieList method on the PG interface.
// Should any errors arise during this retrieval process, it logs the error, responds to the HTTP request with a 502 Bad Gateway status,
// indicating a problem with accessing or processing the data, and returns the error. On successful retrieval, it sends the list of movies
// back to the client in JSON format, providing a comprehensive view of the available movie records.
func (view *View) MovieListView() error {

	log.Info().Msg("MovieListView called")

	data, err := view.PG.MovieList(view.W, view.R)
	if err != nil {
		log.Error().Err(err).Msg("Error in MovieList")
		view.handleError(err, http.StatusBadGateway)
		return err
	}

	view.respondWithJSON(data)
	return nil
}

func (view *View) MovieFindView() error {

	log.Info().Msg("MovieFindView called")

	data, err := view.PG.MovieFind(view.W, view.R)
	if err != nil {
		log.Error().Err(err).Msg("Error in MovieFind")
		view.handleError(err, http.StatusBadGateway)
		return err
	}

	view.respondWithJSON(data)
	return nil
}

// MovieDeleteView oversees the HTTP request for deleting a specific movie.
// The function logs the start of the deletion process, then attempts to delete the specified movie by invoking the MovieDelete method
// on the PG interface. If this deletion process fails, due to reasons like the movie not existing or database constraints,
// it logs the failure, issues a 502 Bad Gateway status in the HTTP response to indicate the inability to process the request,
// and returns the error. If the deletion is successful, it may respond with confirmation in JSON format, indicating the successful removal
// of the movie record.
func (view *View) MovieDeleteView() error {

	log.Info().Msg("MovieDeleteView called")

	data, err := view.PG.MovieDelete(view.W, view.R)
	if err != nil {
		log.Error().Err(err).Msg("Error in MovieDelete")
		view.handleError(err, http.StatusBadGateway)
		return err
	}

	view.respondWithJSON(data)
	return nil
}
