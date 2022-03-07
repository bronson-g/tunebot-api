package main

import (
	"fmt"
	"net/http"

	"github.com/bronson-g/tunebot-api/endpoint"
	"github.com/bronson-g/tunebot-api/model"
	"github.com/gorilla/mux"
)

func main() {
	err := model.Connect()
	if err != nil {
		fmt.Println("Failed to connect to database")
		fmt.Println(err.Error())
		return
	}
	defer model.Disconnect()

	router := mux.NewRouter()

	router.HandleFunc("/user/register/", endpoint.Register).Methods("POST")
	router.HandleFunc("/user/login/", endpoint.Login).Methods("POST")
	router.HandleFunc("/playlist/create/", endpoint.Create).Methods("POST")
	router.HandleFunc("/playlist/update/", endpoint.Update).Methods("POST")
	router.HandleFunc("/playlist/delete/", endpoint.Delete).Methods("POST")
	router.HandleFunc("/playlist/song/add/", endpoint.Add).Methods("POST")
	router.HandleFunc("/playlist/song/remove/", endpoint.Remove).Methods("POST")
	http.ListenAndServe(":8080", router)
}
