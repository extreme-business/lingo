package web

import (
	"fmt"
	"net/http"
)

// Options options for the web server
type Options struct {
	Port int
}

// Start sets up the web server and routes
func Start(opt *Options) error {
	if opt.Port == 0 {
		return fmt.Errorf("port is not set in options")
	}

	server := New()

	http.HandleFunc("/register", server.Register)
	http.HandleFunc("/login", server.Login)

	if err := http.ListenAndServe(fmt.Sprintf(":%d", opt.Port), nil); err != nil {
		return fmt.Errorf("failed to start web server: %w", err)
	}

	return nil
}
