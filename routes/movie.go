package routes

import (
	"net/http"

	"vk.com/m/views"
)

func (router *Router) MovieAddRoute(w http.ResponseWriter, r *http.Request) {
	view := views.View{W: w, R: r, PG: router.PG}
	view.MovieAddView()
}

func (router *Router) MovieEditRoute(w http.ResponseWriter, r *http.Request) {
	view := views.View{W: w, R: r, PG: router.PG}
	view.MovieEditView()
}

func (router *Router) MovieListRoute(w http.ResponseWriter, r *http.Request) {
	view := views.View{W: w, R: r, PG: router.PG}
	view.MovieListView()
}

func (router *Router) MovieFindRoute(w http.ResponseWriter, r *http.Request) {
	view := views.View{W: w, R: r, PG: router.PG}
	view.MovieFindView()
}

func (router *Router) MovieDeleteRoute(w http.ResponseWriter, r *http.Request) {
	view := views.View{W: w, R: r, PG: router.PG}
	view.MovieDeleteView()
}
