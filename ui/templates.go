// Package ui
package ui

import "html/template"

func LoadTemplates() *template.Template {
	t := template.New("")
	template.Must(t.ParseFiles(
		"ui/templates/base.html",
		"ui/templates/index.html",
	))
	// template.Must(t.ParseGlob("templates/partials/*.html"))
	return t
}
