package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/marsorm/goPageDB/pkg/tmpl"
)

type templateMessage struct {
	Message string
}

func NewtemplateMessage(str string) templateMessage {
	return templateMessage{
		Message: str,
	}
}

// APIServer represents your application server.
type APIServer struct {
	port   string
	router *http.ServeMux
}

// NewAPIServer creates a new APIServer and loads templates.
func NewAPIServer(port string) *APIServer {
	mux := http.NewServeMux()

	// Load templates via the separate package.
	if err := tmpl.LoadTemplates("static/*.html"); err != nil {
		log.Fatalf("Could not load templates: %v", err)
	}

	return &APIServer{
		port:   port,
		router: mux,
	}
}

// registerRoutes sets up all HTTP routes for the server.
func (s *APIServer) registerRoutes() {
	s.router.HandleFunc("/", s.landingPage)
	s.router.HandleFunc("/import", s.importHandler)
	s.router.HandleFunc("/export", s.exportHandler)
	s.router.HandleFunc("/help", s.helpPage)
	s.router.HandleFunc("/test", s.test)

	// Serve static assets (like styles.css) from the static directory.
	fs := http.FileServer(http.Dir("static"))
	s.router.Handle("/static/", http.StripPrefix("/static/", fs))
}

// landingPage renders the landing page.
func (s *APIServer) landingPage(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.RenderTemplate(w, "landing", nil); err != nil {
		http.Error(w, "Error rendering landing page", http.StatusInternalServerError)
		log.Printf("landingPage error: %v", err)
	}
}

func (s *APIServer) handleFileAndData(r *http.Request) (string, string, error) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return "", "", err
	}

	var filename string
	fileProvided := false

	if file, handler, err := r.FormFile("datafile"); err == nil {
		defer file.Close()
		if handler.Size == 0 {
			return "", "", fmt.Errorf("uploaded file is empty")
		}
		filename = handler.Filename
		fileProvided = true
	} else {
		log.Printf("No file uploaded or error retrieving file: %v", err)
	}

	pastedData := r.FormValue("data")
	pastedProvided := pastedData != ""

	if fileProvided && pastedProvided {
		return "", "", fmt.Errorf("please provide either an uploaded file OR pasted data, not both")
	}
	if !fileProvided && !pastedProvided {
		return "", "", fmt.Errorf("please provide either an uploaded file or pasted data")
	}

	return filename, pastedData, nil
}

// importHandler serves the import page on GET and processes file/data submissions on POST.
func (s *APIServer) importHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("importHandler called - Method: %s, URL: %s", r.Method, r.URL.String())

	switch r.Method {
	case http.MethodGet:
		if err := tmpl.RenderTemplate(w, "import.html", nil); err != nil {
			http.Error(w, "Error rendering import page", http.StatusInternalServerError)
			log.Printf("importHandler GET error: %v", err)
		}
	case http.MethodPost:
		log.Println("POST request received at /import")
		filename, pastedData, err := s.handleFileAndData(r)
		if err != nil {
			if err := tmpl.RenderTemplate(w, "error", NewtemplateMessage("big error")); err != nil {
				log.Printf("importHandler error rendering error template: %v", err)
			}
			return
		}

		log.Printf("Import - Uploaded file: %s", filename)
		log.Printf("Import - Pasted data: %s", pastedData)

		if err := tmpl.RenderTemplate(w, "success", NewtemplateMessage("Data Imported Successfully!")); err != nil {
			http.Error(w, "Error rendering success page", http.StatusInternalServerError)
			log.Printf("Error rendering success template: %v", err)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("Method not allowed: %s", r.Method)
	}
}

func (s *APIServer) exportHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("exportHandler called - Method: %s, URL: %s", r.Method, r.URL.String())
	// Your export logic here.
}

func (s *APIServer) helpPage(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.RenderTemplate(w, "help", nil); err != nil {
		http.Error(w, "Error rendering help page", http.StatusInternalServerError)
		log.Printf("helpPage error: %v", err)
	}
}

func (s *APIServer) Run() {
	s.registerRoutes()
	log.Printf("Server running on port: %s", s.port)
	if err := http.ListenAndServe(s.port, s.router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func (s *APIServer) test(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.RenderTemplate(w, "index", "woo"); err != nil {
		http.Error(w, "Woo Error rendering template", http.StatusInternalServerError)
	}
}
