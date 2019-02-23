package controllers

import (
	"net/http"

	"github.com/gorilla/schema"
)

//Parses a form and decodes it using the gorilla schema package.
//Used by users.Create to parse the form to be able to pass it to the webapp
func parseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	dec := schema.NewDecoder()
	if err := dec.Decode(dst, r.PostForm); err != nil {
		return err
	}
	return nil
}
