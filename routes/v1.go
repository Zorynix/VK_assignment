package routes

import (
	"net/http"
)

func (router *Router) V1Routes() {
	http.HandleFunc("/v1/actor-add", router.ActorAddRoute)
	http.HandleFunc("/v1/actor-edit/", router.ActorEditRoute)
	http.HandleFunc("/v1/actor-list", router.ActorListRoute)
	http.HandleFunc("/v1/actor-delete/", router.ActorDeleteRoute)
	http.HandleFunc("/v1/movie-add", router.MovieAddRoute)
	http.HandleFunc("/v1/movie-edit/", router.MovieEditRoute)
	http.HandleFunc("/v1/movie-find", router.MovieFindRoute)
	http.HandleFunc("/v1/movie-list", router.MovieListRoute)
	http.HandleFunc("/v1/movie-delete/", router.MovieDeleteRoute)
}
