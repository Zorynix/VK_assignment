package routes

import (
	"net/http"
)

func (router *Router) V1Routes() {
	http.HandleFunc("/v1/", router.)
	http.HandleFunc("/v1/", router.ApiCurrenciesRoute)
}
