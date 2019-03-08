package main

import (
	"fmt"
	"net/http"
	"webapp/controllers"
	"webapp/models"

	"github.com/gorilla/mux"
)

//Constants to connect to the database
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "0102"
	dbname   = "webapp_dev"
)

func main() {
	//Sets up the connection for the database using the constants above.
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	//Connects to the DB, defers the closing of it until application shutdown and calls the Automigrate function.
	//NewUserService, Close and Automigrate are all found in the User models file.
	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()
	us.AutoMigrate()

	//References the functions that load the views. the (us) is the argument passed when the function is called.
	usersC := controllers.NewUsers(us)
	staticC := controllers.NewStatic()

	//Sets up a new router, r.
	//Tells which func to pass to depending on which route is entered
	//ListenandServe sets up the server
	r := mux.NewRouter()

	// Assets
	//Our asset fileserver
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)
	//Static pages are handled in the static controller
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.Handle("/faq", staticC.Faq).Methods("GET")
	r.Handle("/login", usersC.LoginView).Methods("GET")
	//User pages are handled in the user controller
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	http.ListenAndServe(":3000", r)
}
