package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

// APIServer represents your application server.
type APIServer struct {
	port      string
	router    *http.ServeMux
	templates *template.Template
}

// NewAPIServer creates a new APIServer, parsing all templates from the "static" folder.
func NewAPIServer(port string) *APIServer {
	mux := http.NewServeMux()

	tpls, err := template.ParseGlob("static/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	return &APIServer{
		port:      port,
		router:    mux,
		templates: tpls,
	}
}

// registerRoutes sets up all HTTP routes for the server.
func (s *APIServer) registerRoutes() {
	s.router.HandleFunc("/", s.landingPage)
	s.router.HandleFunc("/import", s.importHandler)
	s.router.HandleFunc("/export", s.exportHandler)
	s.router.HandleFunc("/help", s.helpPage)

	// Serve static assets (like styles.css) from the static directory.
	fs := http.FileServer(http.Dir("static"))
	s.router.Handle("/static/", http.StripPrefix("/static/", fs))
}

// landingPage renders the landing page.
func (s *APIServer) landingPage(w http.ResponseWriter, r *http.Request) {
	if err := s.templates.ExecuteTemplate(w, "landing.html", nil); err != nil {
		http.Error(w, "Error rendering landing page", http.StatusInternalServerError)
		log.Printf("landingPage error: %v", err)
	}
}

func (s *APIServer) handleFileAndData(r *http.Request) (string, string, error) {
	// Parse the multipart form with a 10MB limit.
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

	// Enforce that only one is provided.
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
		if err := s.templates.ExecuteTemplate(w, "import.html", nil); err != nil {
			http.Error(w, "Error rendering import page", http.StatusInternalServerError)
			log.Printf("importHandler GET error: %v", err)
		}
	case http.MethodPost:
		log.Println("POST request received at /import")
		filename, pastedData, err := s.handleFileAndData(r)
		if err != nil {
			http.Error(w, "Error processing form data", http.StatusInternalServerError)
			log.Printf("Error in handleFileAndData: %v", err)
			return
		}

		//Debug print
		log.Printf("Import - Uploaded file: %s", filename)
		log.Printf("Import - Pasted data: %s", pastedData)

		data := struct {
			Message string
		}{
			Message: "Data Imported Successfully!",
		}
		if err := s.templates.ExecuteTemplate(w, "success.html", data); err != nil {
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

	switch r.Method {
	case http.MethodGet:
		if err := s.templates.ExecuteTemplate(w, "export.html", nil); err != nil {
			http.Error(w, "Error rendering export page", http.StatusInternalServerError)
			log.Printf("exportHandler GET error: %v", err)
		}
	case http.MethodPost:
		log.Println("POST request received at /export")
		filename, pastedData, err := s.handleFileAndData(r)
		if err != nil {
			http.Error(w, "Error processing form data", http.StatusInternalServerError)
			log.Printf("Error in handleFileAndData: %v", err)
			return
		}

		log.Printf("Export - Uploaded file: %s", filename)
		log.Printf("Export - Pasted data: %s", pastedData)

		// Render the success template.
		data := struct {
			Message string
		}{
			Message: "Data Exported Successfully!",
		}
		if err := s.templates.ExecuteTemplate(w, "success.html", data); err != nil {
			http.Error(w, "Error rendering success page", http.StatusInternalServerError)
			log.Printf("Error rendering success template: %v", err)
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("Method not allowed: %s", r.Method)
	}
}

// helpPage renders the help page.
func (s *APIServer) helpPage(w http.ResponseWriter, r *http.Request) {
	if err := s.templates.ExecuteTemplate(w, "help.html", nil); err != nil {
		http.Error(w, "Error rendering help page", http.StatusInternalServerError)
		log.Printf("helpPage error: %v", err)
	}
}

// Run starts the HTTP server.
func (s *APIServer) Run() {
	s.registerRoutes()
	log.Printf("Server running on port: %s", s.port)
	if err := http.ListenAndServe(s.port, s.router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
