package views

import (
	_ "embed"
	"text/template"
)

// Embed the layout HTML template
//
//go:embed layout.html
var layoutHTML string

// Parse the layout and error templates
var layoutTemplate = template.Must(template.New("layout").Parse(layoutHTML))
