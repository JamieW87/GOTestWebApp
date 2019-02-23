package controllers

import "webapp/views"

type Static struct {
	Home    *views.View
	Contact *views.View
	Faq     *views.View
}

//NewStatic Initializes the views for the static pages. Instead of doing it in main.go
//Sends necessary information to NewView
func NewStatic() *Static {
	return &Static{
		Home: views.NewView(
			"bootstrap", "static/home"),
		Contact: views.NewView(
			"bootstrap", "static/contact"),
		Faq: views.NewView(
			"bootstrap", "static/faq"),
	}
}
