package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/EleisonC/User-API.git/config"
	"github.com/EleisonC/User-API.git/routes"
)

func main() {
	r := mux.NewRouter()
	routes.RegDogOwnerRoutes(r)
	http.Handle("/", r)
	config.ConnectDB()
	log.Fatal(http.ListenAndServe(":8081", r))
}