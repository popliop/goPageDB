package server

import (
	"fmt"
	"io"
	"net/http"
)

type APIServer struct {
	serverPort string
	router     http.ServeMux
}

func (s *APIServer) Run() {
	s.registerRoutes()

	fmt.Println("Server running on port: ", s.serverPort)
	http.ListenAndServe(s.serverPort, &s.router)
}

func NewAPIServer(port string) *APIServer {
	return &APIServer{
		serverPort: port,
		router:     *http.NewServeMux(),
	}
}

func (s *APIServer) registerRoutes() {
	s.router.HandleFunc("/", s.homepage)
	s.router.HandleFunc("/submit", s.Form)
	s.router.HandleFunc("/styles.css", s.styles)
}

func (s *APIServer) homepage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
	fmt.Println(r.URL.Path)
}

func (s *APIServer) Form(w http.ResponseWriter, r *http.Request) {
	req, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	r.Body.Close()
	fmt.Println(string(req))

	response := "Tack"
	rpB := []byte(response)

	w.Write(rpB)

}

func (s *APIServer) styles(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/styles.css")
	fmt.Println(r.URL.Path)
}
