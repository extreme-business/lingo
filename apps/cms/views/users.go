package views

import (
	"io"
	"text/template"

	"github.com/extreme-business/lingo/apps/account/domain"

	_ "embed"
)

//go:embed userlist.html
var userListHTML string
var userListTemplate = template.Must(template.Must(layoutTemplate.Clone()).Parse(userListHTML))

func UserList(w io.Writer, users []*domain.User) error {
	data := map[string]interface{}{
		"Users": users,
	}
	return userListTemplate.Execute(w, data)
}
