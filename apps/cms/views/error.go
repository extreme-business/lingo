package views

import (
	"io"
	"text/template"
	"time"

	_ "embed"
)

//go:embed error.html
var errorHTML string
var errorTemplate = template.Must(template.Must(layoutTemplate.Clone()).Parse(errorHTML))

func Error(w io.Writer, message string) error {
	// Create a map to hold the data, including the current year
	data := map[string]interface{}{
		"Message": message,
		"Year":    time.Now().Year(),
	}

	return errorTemplate.Execute(w, data)
}
