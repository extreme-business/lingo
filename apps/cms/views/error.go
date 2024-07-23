package views

import (
	"io"
	"text/template"

	_ "embed"
)

//go:embed error.html
var errorHTML string
var errorTemplate = template.Must(template.Must(layoutTemplate.Clone()).Parse(errorHTML))

func ShowError(w io.Writer, message string) error {
	return errorTemplate.Execute(w, message)
}
