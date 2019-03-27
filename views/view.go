package views

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
)

//Declaring global variables
//LayoutDir specifies the layout directory and TemplateExt states the extension we want our files to match
var (
	LayoutDir   = "views/layouts/"
	TemplateDir = "views/"
	TemplateExt = ".gohtml"
)

//View declares the View type
type View struct {
	Template *template.Template
	Layout   string
}

//NewView runs logic for all views
/*This function handles appending common template files to the list of files provided*/
//Func main calls this function, passes to it the layout (a string) and the files it must load.
func NewView(layout string, files ...string) *View {
	//Uses the two functions below to append and prepend file paths and extensions.
	addTemplatePath(files)
	addTemplateExt(files)
	//parses the files in the layout folder as set out in layoutFiles function
	files = append(files, layoutFiles()...)
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	//This function then returns the parsed templates, and the layout that was passed to it.
	return &View{
		Template: t,
		Layout:   layout,
	}
}

//Uses the global variables above and returns a slice of a string.
//Basically a shortener for the layout files.
func layoutFiles() []string {
	//States the layouts we are going to include in our views (Anything in the layoutDir) with the filename all (*) with
	//the extension that matches TemplateExt... Basically creates the variable files which equals "views/layouts/*.gohtml"
	//Globs them together so that we can clean up our code and not have to reference each individual one.
	//Purely for clean up purposes.
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	//Returns the files variable set above, as a slice of string.
	return files
}

//These two functions prepend the file paths and append the extensions for files that are passed into NewView. Simpliyfying our code.
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}

func (v *View) Render(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	switch data.(type) {
	case Data:
		// do nothing
	default:
		data = Data{
			Content: data,
		}
	}
	var buf bytes.Buffer
	err := v.Template.ExecuteTemplate(&buf, v.Layout, data)
	if err != nil {
		http.Error(w, "Something went wrong. If the problem "+
			"persists, please email support@lenslocked.com",
			http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

//Used to render the pages
func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, nil)
}
