package views

import (
	"html/template"
	"net/http"
	"path/filepath"
)

//Declaring global variables
//LayoutDir specifies the layout directory and TemplateExt states the extension we want our files to match
var (
	LayoutDir   = "views/layouts/"
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
func layoutFiles() []string {
	//States the templates we are going to include in our view, Anything in the layoutDir, with the filename all (*) with
	//the extension that matches TemplateExt = Basically creates the variable files which equals "views/layouts/*.gohtml"
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		panic(err)
	}
	//Returns the files variable set above, as a slice of string.
	return files
}

//Render method added to the View type (Why the view part comes before the function name)
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	return v.Template.ExecuteTemplate(w, v.Layout, data)
}
