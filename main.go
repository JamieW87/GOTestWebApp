package main

import (
	"net/http"
	"webapp/views"

	"github.com/gorilla/mux"
)

//Global variables
var (
	homeView    *views.View
	contactView *views.View
	faqView     *views.View
	signupView  *views.View
)

func main() {
	//Builds our views when the program first starts running.
	//Calls the NewView function, passes it the bootstrap template and the page it wants it to open.
	homeView = views.NewView("bootstrap", "views/home.gohtml")
	contactView = views.NewView("bootstrap", "views/contact.gohtml")
	faqView = views.NewView("bootstrap", "views/faq.gohtml")
	signupView = views.NewView("bootstrap", "views/signup.gohtml")

	//Handlefuncs tell which route starts which func.
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/faq", faq)
	r.HandleFunc("/signup", signup)
	http.ListenAndServe(":3000", r)
}

/*Below are the functions that deal with the requests from the browser.*/
func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	//So now the Render function deals with executing the template to keep the functions clean
	must(homeView.Render(w, nil))
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(contactView.Render(w, nil))
}

func faq(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(faqView.Render(w, nil))
}

func signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	must(signupView.Render(w, nil))
}

// A helper function that panics on any error
func must(err error) {
	if err != nil {
		panic(err)
	}
}
