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

	//TODO: change the http methods to reflect behaviour, DELETE/PUT etc not all POST
	router.HandleFunc("/user/register/", endpoint.Register).Methods("POST")
	router.HandleFunc("/user/login/", endpoint.Login).Methods("POST")
	router.HandleFunc("/device/user/link/", endpoint.Link).Methods("POST")
	router.HandleFunc("/device/user/get/", endpoint.Get).Methods("POST")
	router.HandleFunc("/playlist/create/", endpoint.Create).Methods("POST")
	router.HandleFunc("/playlist/update/", endpoint.Update).Methods("POST")
	router.HandleFunc("/playlist/delete/", endpoint.Delete).Methods("POST")
	router.HandleFunc("/playlist/song/add/", endpoint.Add).Methods("POST")
	router.HandleFunc("/playlist/song/remove/", endpoint.Remove).Methods("POST")

	log.Println(log.Green("Listening on port 8080."))
	http.ListenAndServe(":8080", router)
}
