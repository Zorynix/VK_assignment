package routes

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "vk.com/m/docs"

	"vk.com/m/middleware"
)

func (router *Router) V1Routes() {

	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	http.HandleFunc("/v1/auth", router.LoginHandler)

	http.Handle("/v1/actor-add", middleware.AuthMiddleware(http.HandlerFunc(router.ActorAddRoute), "admin"))
	http.Handle("/v1/actor-edit/", middleware.AuthMiddleware(http.HandlerFunc(router.ActorEditRoute), "admin"))
	http.Handle("/v1/actor-list", middleware.AuthMiddleware(http.HandlerFunc(router.ActorListRoute), "admin", "user"))
	http.Handle("/v1/actor-delete/", middleware.AuthMiddleware(http.HandlerFunc(router.ActorDeleteRoute), "admin"))

	http.Handle("/v1/movie-add", middleware.AuthMiddleware(http.HandlerFunc(router.MovieAddRoute), "admin"))
	http.Handle("/v1/movie-edit/", middleware.AuthMiddleware(http.HandlerFunc(router.MovieEditRoute), "admin"))
	http.Handle("/v1/movie-list", middleware.AuthMiddleware(http.HandlerFunc(router.MovieListRoute), "admin", "user"))
	http.Handle("/v1/movie-find", middleware.AuthMiddleware(http.HandlerFunc(router.MovieFindRoute), "admin", "user"))
	http.Handle("/v1/movie-delete/", middleware.AuthMiddleware(http.HandlerFunc(router.MovieDeleteRoute), "admin"))
}
