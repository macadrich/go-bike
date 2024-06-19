package main

import (
	"log"
	"net/http"

	"github.com/macadrich/go-bike/api"
	"github.com/macadrich/go-bike/api/handlers"
	"github.com/macadrich/go-bike/api/routers"
	"github.com/macadrich/go-bike/client"
	"github.com/macadrich/go-bike/config"

	"github.com/macadrich/go-bike/database/postgres"
)

func main() {

	db, err := postgres.NewDB(config.LoadDBConfig())
	if err != nil {
		log.Fatal(err)
	}

	client := client.NewClient()
	service := api.NewService(db, client)
	handlers := handlers.NewHandlers(service)
	router := routers.NewRouter(handlers)

	log.Println("Server is running on: 8080")
	http.ListenAndServe(":8080", router)
}
