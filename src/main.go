package main

import (
	"net/http"

	"github.com/bronson-g/tunebot-api/endpoint"
	"github.com/bronson-g/tunebot-api/log"
	"github.com/bronson-g/tunebot-api/model"
	"github.com/gorilla/mux"
)

func main() {
	log.Println("tunebot-api")

	err := model.Connect()
	if err != nil {
		log.Println(log.Red("Failed to connect to database."))
		log.Println(log.Red(err.Error()))
		return
	} else {
		log.Println(log.Green("Connected to database."))
	}
	defer model.Disconnect()

	router := mux.NewRouter()

	router.HandleFunc("/user/register/", endpoint.Register).Methods("POST")
	router.HandleFunc("/user/login/", endpoint.Login).Methods("POST")
	router.HandleFunc("/device/user/link/", endpoint.Link).Methods("LINK")
	router.HandleFunc("/device/user/get/", endpoint.Get).Methods("GET")
	router.HandleFunc("/playlist/create/", endpoint.Create).Methods("POST")
	router.HandleFunc("/playlist/update/", endpoint.Update).Methods("PATCH")
	router.HandleFunc("/playlist/delete/", endpoint.Delete).Methods("DELETE")
	router.HandleFunc("/playlist/song/add/", endpoint.Add).Methods("PUT")
	router.HandleFunc("/playlist/song/remove/", endpoint.Remove).Methods("DELETE")

	log.Println(log.Green("Listening on port 8080."))
	http.ListenAndServe(":8080", router)
}
