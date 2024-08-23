package views

import (
	"io"
	"text/template"
	"time"

	"github.com/extreme-business/lingo/apps/account/domain"

	_ "embed"
)

var (
	//go:embed layout.html
	layoutHTML string
	//go:embed login.html
	loginHTML string
	//go:embed logout.html
	logoutHTML string
	//go:embed error.html
	errorHTML string
	//go:embed userlist.html
	userListHTML string
)

type Writer struct {
	layoutTemplate   *template.Template
	loginTemplate    *template.Template
	logoutTemplate   *template.Template
	errorTemplate    *template.Template
	userListTemplate *template.Template
}

func New() (*Writer, error) {
	tmpl := template.Must(template.New("layout").Parse(layoutHTML))

	// Create sub-templates by cloning and parsing each specific template
	loginTemplate, err := template.New("login").Parse(loginHTML)
	if err != nil {
		return nil, err
	}

	logoutTemplate, err := template.New("logout").Parse(loginHTML)
	if err != nil {
		return nil, err
	}

	errorTemplate := template.Must(tmpl.Clone()).New("error")
	template.Must(errorTemplate.Parse(errorHTML))

	userListTemplate := template.Must(tmpl.Clone()).New("userlist")
	template.Must(userListTemplate.Parse(userListHTML))

	return &Writer{
		layoutTemplate:   tmpl,
		loginTemplate:    loginTemplate,
		logoutTemplate:   logoutTemplate,
		errorTemplate:    errorTemplate,
		userListTemplate: userListTemplate,
	}, nil
}

func (v *Writer) Login(w io.Writer) error {
	data := map[string]interface{}{
		"Year": time.Now().Year(),
	}

	return v.loginTemplate.ExecuteTemplate(w, "layout", data)
}

func (v *Writer) Logout(w io.Writer) error {
	data := map[string]interface{}{
		"Year": time.Now().Year(),
	}

	return v.logoutTemplate.ExecuteTemplate(w, "layout", data)
}

func (v *Writer) Error(w io.Writer, message string) error {
	data := map[string]interface{}{
		"Message": message,
		"Year":    time.Now().Year(),
	}

	return v.errorTemplate.ExecuteTemplate(w, "layout", data)
}

func (v *Writer) UserList(w io.Writer, users []*domain.User) error {
	data := map[string]interface{}{
		"Users": users,
		"Year":  time.Now().Year(),
	}

	return v.userListTemplate.ExecuteTemplate(w, "layout", data)
}
