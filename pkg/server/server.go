package server

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/marsorm/goPageDB/pkg/tmpl"
)

type templateMessage struct {
	Message string
}

func NewTemplateMessage(str string) templateMessage {
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

func validateForm(header *multipart.FileHeader) error {
	if header.Size == 0 {
		return fmt.Errorf("uploaded file is empty")
	}

	if filepath.Ext(header.Filename) != ".csv" {
		return fmt.Errorf("invalid extention")
	}
	fmt.Println(header)
	return nil
}

func parseRequestForm(r *http.Request) (multipart.File, *multipart.FileHeader, error) {
	fmt.Println("hello from the othereside")
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, nil, err
	}

	file, header, err := r.FormFile("datafile")
	if err != nil {
		return nil, nil, fmt.Errorf("unable to retrieve form file from request %w", err)
	}
	return file, header, nil
}

func (s *APIServer) handleFileAndData(r *http.Request) error {
	fmt.Println("!!!!!!")
	_, header, err := parseRequestForm(r)
	if err != nil {
		return fmt.Errorf("handleFileAndData: %w", err)
	}

	if err := validateForm(header); err != nil {
		return fmt.Errorf("validateform: %w", err)
	}

	return nil
}

// importHandler serves the import page on GET and processes file/data submissions on POST.
func (s *APIServer) importHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("importHandler called - Method: %s, URL: %s", r.Method, r.URL.String())

	switch r.Method {
	case http.MethodGet:
		// Render the import page for GET requests.
		if err := tmpl.RenderTemplate(w, "import.html", nil); err != nil {
			log.Printf("Error rendering import page: %v", err)
			http.Error(w, "Error rendering import page", http.StatusInternalServerError)
		}

	case http.MethodPost:
		log.Println("POST request received at /import")
		// Process the incoming file and data.
		if err := s.handleFileAndData(r); err != nil {
			log.Printf("Error handling file and data: %v", err)
			// Try to render an error template with the error message.
			if renderErr := tmpl.RenderTemplate(w, "error", NewTemplateMessage(err.Error())); renderErr != nil {
				log.Printf("Error rendering error template: %v", renderErr)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		// If everything went well, render a success page.
		if err := tmpl.RenderTemplate(w, "success", NewTemplateMessage("Data Imported Successfully!")); err != nil {
			log.Printf("Error rendering success template: %v", err)
			http.Error(w, "Error rendering success page", http.StatusInternalServerError)
		}

	default:
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
