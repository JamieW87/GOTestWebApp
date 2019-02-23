package main

import (
	"net/http"
	"webapp/controllers"

	"github.com/gorilla/mux"
)

func main() {
	//References the functions that load the views
	usersC := controllers.NewUsers()
	staticC := controllers.NewStatic()

	//Sets up a new router, r.
	//Tells which func to pass to depending on which route is entered
	//ListenandServe sets up the server
	r := mux.NewRouter()
	//Static pages are handled in the static controller
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.Faq).Methods("GET")
	//User pages are handled in the user controller
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	http.ListenAndServe(":3000", r)
}
