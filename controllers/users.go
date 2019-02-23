package controllers

import (
	"fmt"
	"net/http"
	"webapp/views"
)

type Users struct {
	NewView *views.View
}

//The struct tags on the end are setting up the input fields for the schema package.
type SignupForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

//NewUsers Handles the loading of the view
func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "users/new"),
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
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	//Uses the parseForm function in helpers to parse and decode the form to be able to pass it to the webapp
	var form SignupForm
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	fmt.Fprintln(w, "Email is", form.Email)
	fmt.Fprintln(w, "Password is", form.Password)
}
