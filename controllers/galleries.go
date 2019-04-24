package controllers

import (
	"net/http"
	"strconv"
	"fmt"

	"webapp/context"
	"webapp/models"
	"webapp/views"

	"github.com/gorilla/mux"
)

type Galleries struct {
	New      *views.View
	ShowView *views.View
	EditView *views.View
	gs       models.GalleryService
	r        *mux.Router
}

type GalleryForm struct {
	Title string `schema:"title"`
}

const (
	ShowGallery = "show_gallery"
)

//NewGalleries passes in a galleryService to use, and deals with the handler function in main.go
func NewGalleries(gs models.GalleryService, r *mux.Router) *Galleries {
	return &Galleries{
		New:      views.NewView("bootstrap", "galleries/new"),
		ShowView: views.NewView("bootstrap", "galleries/show"),
		EditView: views.NewView("bootstrap", "galleries/edit"),
		gs:       gs,
		r:        r,
	}
}

// POST /galleries
func (g *Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var vd views.Data
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	user := context.User(r.Context())
	gallery := models.Gallery{
		Title:  form.Title,
		UserID: user.ID,
	}
	if err := g.gs.Create(&gallery); err != nil {
		vd.SetAlert(err)
		g.New.Render(w, vd)
		return
	}
	url, err := g.r.Get(ShowGallery).URL("id",
		strconv.Itoa(int(gallery.ID)))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	http.Redirect(w, r, url.Path, http.StatusFound)
}

// POST /galleries/:id/update
func (g *Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "Gallery not found", http.StatusNotFound)
		return
	}
	var vd views.Data
	vd.Content = gallery
	var form GalleryForm
	if err := parseForm(r, &form); err != nil {
		// If there is an error we are going to render the
		// EditView again with an alert message.
		vd.SetAlert(err)
		g.EditView.Render(w, vd)
		return
	}
	gallery.Title = form.Title
	err = g.gs.Update(gallery)
	//The if else renders an alert if theres an error, if not it still renders an alert
	//but as a success message.
	if err != nil {
		vd.SetAlert(err)
	} else {
		vd.Alert = &views.Alert {
			Level: views.AlertLvlSuccess,
			Message: "Gallery Succesfully updated",
		}
		g.EditView.Render(w, vd)
	}
	
	vd.Alert = &views.Alert{
		Level:   views.AlertLvlSuccess,
		Message: "Gallery updated successfully!",
	}
	g.EditView.Render(w, vd)

}

// GET /galleries/:id
func (g *Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	var vd views.Data
	vd.Content = gallery
	g.ShowView.Render(w, vd)
}

//GET /galleries/:id/edit
func (g *Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You do not have permission to edit this gallery", http.StatusForbidden)
		return
	}
	var vd views.Data
	vd.Content = gallery
	g.EditView.Render(w, vd)
}

// POST /galleries/:id/delete
func (g *Galleries) Delete(w http.ResponseWriter, r *http.Request) {
	// Lookup the gallery using the galleryByID we wrote earlier
	gallery, err := g.galleryByID(w, r)
	if err != nil {
	// If there is an error the galleryByID will have rendered
	// it for us already.
	return
	}
	// We also need to retrieve the user and verify they have
	// permission to delete this gallery. This means we will
	// need to use the RequireUser middleware on any routes
	// mapped to this method.
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
	http.Error(w, "You do not have permission to edit "+
	"this gallery", http.StatusForbidden)
	return
	}
	var vd views.Data
	err = g.gs.Delete(gallery.ID)
	if err != nil {
	// If there is an error we want to set an alert and
	// render the edit page with the error. We also need
	// to set the Yield to gallery so that the EditView
	// is rendered correctly.
	vd.SetAlert(err)
	vd.Yield = gallery
	g.EditView.Render(w, vd)
	return
	}
	// TODO: We will eventually want to redirect to the index
	// page that lists all galleries this user owns, but for
	// now a success message will suffice.
	fmt.Fprintln(w, "successfully deleted!")
	}


func (g *Galleries) galleryByID(w http.ResponseWriter,
	r *http.Request) (*models.Gallery, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusNotFound)
		return nil, err
	}
	gallery, err := g.gs.ByID(uint(id))
	if err != nil {
		switch err {
		case models.ErrNotFound:
			http.Error(w, "Gallery not found", http.StatusNotFound)
		default:
			http.Error(w, "Whoops! Something went wrong.", http.StatusInternalServerError)
		}
		return nil, err
	}
	return gallery, nil
}
