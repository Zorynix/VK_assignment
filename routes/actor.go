package routes

import (
	"net/http"

	"vk.com/m/views"
)

func (router *Router) ActorAddRoute(w http.ResponseWriter, r *http.Request) {
	view := views.View{W: w, R: r, PG: router.PG}
	view.ActorAddView()
}

func (router *Router) ActorEditRoute(w http.ResponseWriter, r *http.Request) {
	view := views.View{W: w, R: r, PG: router.PG}
	view.ActorEditView()
}

func (router *Router) ActorListRoute(w http.ResponseWriter, r *http.Request) {
	view := views.View{W: w, R: r, PG: router.PG}
	view.ActorListView()
}

func (router *Router) ActorDeleteRoute(w http.ResponseWriter, r *http.Request) {
	view := views.View{W: w, R: r, PG: router.PG}
	view.ActorDeleteView()
}
