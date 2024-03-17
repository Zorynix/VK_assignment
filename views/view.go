package views

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"vk.com/m/services"
)

type View struct {
	W  http.ResponseWriter
	R  *http.Request
	PG *services.Postgresql
}

// respondWithJSON takes any data interface{}, serializes it to JSON, and writes it to the HTTP response.
// It sets the Content-Type header to "application/json" to indicate the MIME type of the response.
// This method is used to send structured data (like objects or arrays) back to the client in a format that's easily parsed and used in web applications.
// In case of serialization failure (e.g., if the data contains unserializable references), it logs the error
// and calls handleError to respond with an HTTP 500 Internal Server Error, indicating a problem with generating the JSON response.
//
// Parameters:
// - data interface{}: The data to serialize into JSON and write to the response.
func (view *View) respondWithJSON(data interface{}) {
	view.W.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(view.W).Encode(data); err != nil {
		view.handleError(err, http.StatusInternalServerError)
	}
}

// handleError logs the provided error and responds to the HTTP request with the specified status code.
// This method standardizes error handling across view functions, ensuring that all errors are logged for debugging purposes
// and that the client receives a consistent error response format. It uses the http.Error utility function
// to send a plain text response containing the standard HTTP status message corresponding to the statusCode.
// This method simplifies the process of returning errors, abstracting the common tasks of setting the response status code
// and logging the error.
//
// Parameters:
// - err error: The error encountered during the handling of the request.
// - statusCode int: The HTTP status code to respond with, indicating the nature of the error (e.g., 400 Bad Request, 404 Not Found).
func (view *View) handleError(err error, statusCode int) {
	log.Info().Err(err).Msg("")
	http.Error(view.W, http.StatusText(statusCode), statusCode)
}
