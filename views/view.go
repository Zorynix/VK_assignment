package views

import (
	"net/http"

	"vk.com/m/services"
)

type View struct {
	W  http.ResponseWriter
	R  *http.Request
	PG *services.Postgresql
}
