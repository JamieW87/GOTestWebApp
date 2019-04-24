package main

import (
	"fmt"
	"net/http"
	"webapp/controllers"
	"webapp/middleware"
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
	services, err := models.NewServices(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer services.Close()
	services.AutoMigrate()
	//services.AutoMigrate()

	//Sets up a new router, r.
	//Tells which func to pass to depending on which route is entered
	r := mux.NewRouter()

	//References the functions that load the views. the (us) is the argument passed when the function is called.
	usersC := controllers.NewUsers(services.User)
	staticC := controllers.NewStatic()
	galleriesC := controllers.NewGalleries(services.Gallery, r)

	//Code for middleware that requires a user to be logged in before they visit certain pages.
	requireUserMw := middleware.RequireUser{
		UserService: services.User,
	}
	// galleriesC.New is an http.Handler, so we use Apply
	newGallery := requireUserMw.Apply(galleriesC.New)
	// galleriecsC.Create is an http.HandlerFunc, so we use ApplyFn
	createGallery := requireUserMw.ApplyFn(galleriesC.Create)

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
	r.HandleFunc("/cookietest", usersC.CookieTest).Methods("GET")
	//Gallery pages - Goes through RequiredUser middleware
	r.Handle("/galleries/new", newGallery).Methods("GET")
	r.Handle("/galleries", createGallery).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name(controllers.ShowGallery)
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET")
	r.HandleFunc("/galleries/id{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")
	//User pages are handled in the user controller
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	http.ListenAndServe(":3000", r)
}
