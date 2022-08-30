package routes


import (
	"github.com/gorilla/mux"
	"github.com/EleisonC/User-API.git/controllers"
)

var RegDogOwnerRoutes = func(router *mux.Router) {
	router.HandleFunc("/createNewDogOwner", controllers.CreateDogOwner).Methods("POST")
	router.HandleFunc("/getAllDogOwners", controllers.GetAllDogOwners).Methods("GET")
	router.HandleFunc("/getDogOwner/{ownerID}", controllers.GetDogOwnerById).Methods("GET")
	router.HandleFunc("/updateOwnerInfo/{ownerID}", controllers.UpdateDogOwner).Methods("PUT")
	router.HandleFunc("/deleteOwner/{ownerID}", controllers.DeleteDogOwner).Methods("DELETE")
}