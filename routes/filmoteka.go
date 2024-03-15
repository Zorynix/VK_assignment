package routes

import (
	"currency-conversion/views"
	"net/http"
)

func (router *Router) ApiExchangeRateRoute(w http.ResponseWriter, r *http.Request) {
	view := views.View{W: w, PG: router.PG}
	view.ExchangeRateView()
}

func (router *Router) ApiCurrenciesRoute(w http.ResponseWriter, r *http.Request) {
	view := views.View{W: w, PG: router.PG}
	view.CurrenciesView()
}

func (router *Router) ApiUpdateRates(w http.ResponseWriter, r *http.Request) {
	view := views.View{W: w, PG: router.PG}
	view.RateHistoryView()
}

func (router *Router) ApiGetRateByCode(w http.ResponseWriter, r *http.Request) {
	view := views.View{W: w, PG: router.PG, R: r}
	view.RateByCodeView()
}
