package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

const (
	what            = "brighella"
	resolverAddress = "8.8.8.8"
	resolverPort    = "53"
)

var (
	httpPort string
)

func init() {
	httpPort = os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "5000"
	}
}

func main() {
	log.Printf("Starting %s...\n", what)

	server := NewServer()

	log.Printf("%s listening on %s...\n", what, httpPort)
	if err := http.ListenAndServe(":"+httpPort, server); err != nil {
		log.Panic(err)
	}
}

// Server represents a front-end web server.
type Server struct {
	// Router which handles incoming requests
	mux *http.ServeMux
}

// NewServer returns a new front-end web server that handles HTTP requests for the app.
func NewServer() *Server {
	router := http.NewServeMux()
	server := &Server{mux: router}
	router.HandleFunc("/", server.Root)
	return server
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// Root is the handler for the HTTP requests to /.
// Any request that is not for the root path / is automatically redirected
// to the root with a 302 status code. Only a request to / will enable the iframe.
func (s *Server) Root(w http.ResponseWriter, r *http.Request) {
	targetURL := os.Getenv("TARGET_URL")
	s.MaskedRedirect(w, r, targetURL)
}

func (s *Server) TemporaryRedirect(w http.ResponseWriter, r *http.Request, strURL string) {
	http.Redirect(w, r, strURL, http.StatusTemporaryRedirect)
}

func (s *Server) MaskedRedirect(w http.ResponseWriter, r *http.Request, strURL string) {
	w.Header().Set("Content-type", "text/html")
	t, _ := template.ParseFiles("redirect.tmpl")
	t.Execute(w, &frame{Src: strURL})
}

type frame struct {
	Src string
}
