package routes

import (
	"context"
	"log"
	"net/http"

	"vk.com/m/services"
)

type Router struct {
	PG *services.Postgresql
}

func Routes(addr *string) {
	postgres, err := services.NewPostgreSQL(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	router := Router{PG: postgres}

	router.V1Routes()

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("Can not start http server: ", err)
	}
}
