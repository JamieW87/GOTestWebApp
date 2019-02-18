package views

import "html/template"

//NewView runs logic for all views
/*This function handles appending common template files to the list of files provided*/
func NewView(layout string, files ...string) *View {
	files = append(files,
		"views/layouts/footer.gohtml",
		"views/layouts/bootstrap.gohtml")
	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}
	return &View{
		Template: t,
		Layout:   layout,
	}
}

//View declares the View type
type View struct {
	Template *template.Template
	Layout   string
}
