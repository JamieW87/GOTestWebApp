package controllers

import (
	"fmt"
	"net/http"
	"webapp/models"
	"webapp/views"
)

//Users adds the view and the UserService to the Users type so that all methods can access these when they need it.
type Users struct {
	NewView *views.View
	us      *models.UserService
}

//The struct tags on the end are setting up the input fields for the schema package.
//These values are taken from the HTML form.
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

//NewUsers Handles the loading of the view and assigns the UserService to
func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
		us:      us,
	}
}

//New Method, or handler, that deals with web requests when a user visits the sign up page
//Uses u to access the Users controller, enables us to reference the NewView field.
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	if err := u.NewView.Render(w, nil); err != nil {
		panic(err)
	}
}

//Create method for POST /signup
//Creates a models.User type via the sign up form
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	//Uses the parseForm function in helpers.go to parse and decode the form to be able to pass it to the webapp
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	//These values are parsed from the signup form and added to the user model.
	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "User is", user)
}
