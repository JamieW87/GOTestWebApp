package controllers

import (
	"fmt"
	"net/http"
	"webapp/models"
	"webapp/rand"
	"webapp/views"
)

//Users adds the view and the UserService to the Users type so that all methods can access these when they need it.
type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        models.UserService
}

//SignupForm holds values for the fields off the form.
//The struct tags on the end are setting up the input fields for the schema package.
//These values are taken from the HTML form.
type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

//LoginForm holds the values from the login.gohtml
type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

//NewUsers Handles the loading of the views and assigns the UserService to
func NewUsers(us models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
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
	err := u.signIn(w, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// Login is used to process the login form when a user
// tries to log in as an existing user (via email & pw).
//
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	form := LoginForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	user, err := u.us.Authenticate(form.Email, form.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			fmt.Fprintln(w, "Invalid email address.")
		case models.ErrInvalidPassword:
			fmt.Fprintln(w, "Invalid password provided.")
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	err = u.signIn(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//Signs in a user via cookies.
//Checks if a raw token is set, if one is not set it creates a new one and updates the user.
func (u *Users) signIn(w http.ResponseWriter, user *models.User) error {
	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
		err = u.us.Update(user)
		if err != nil {
			return err
		}
	}
	cookie := http.Cookie{
		Name:     "remember_token",
		Value:    user.Remember,
		HttpOnly: true,
		//HttpOnly means our cookies arent vulnerable to Javascript attacks. May need disabling to run JS in future.
	}
	http.SetCookie(w, &cookie)
	return nil
}

//CookieTest tests the succesful writing of cookies to the users pc
func (u *Users) CookieTest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("remember_token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user, err := u.us.ByRemember(cookie.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)
}
