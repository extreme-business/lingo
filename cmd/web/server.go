package web

import (
	"fmt"
	"net/http"
)

type Server struct{}

// serveRegisterPage serves the registration page
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	// You can render a registration form here
	if r.Method == http.MethodPost {
		// You can handle the registration here
		fmt.Fprintf(w, "<h1>Registration successful</h1>")
		return
	}

	fmt.Fprintf(w, "<h1>Register</h1><form>Registration form here...</form>")
}

// serveLoginPage serves the login page
func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	// You can render a login form here
	fmt.Fprintf(w, "<h1>Login</h1><form>Login form here...</form>")
}

func New() *Server {
	return &Server{}
}
