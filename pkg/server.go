package server

import (
	"html/template"
	"log"
	"net/http"
)

// APIServer represents your application server.
type APIServer struct {
	serverPort string
	router     *http.ServeMux
	templates  *template.Template
}

// NewAPIServer initializes the server and parses the templates.
func NewAPIServer(port string) *APIServer {
	mux := http.NewServeMux()
	// Parse all HTML templates in the static directory.
	tpls, err := template.ParseGlob("static/*.html")
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}
	return &APIServer{
		serverPort: port,
		router:     mux,
		templates:  tpls,
	}
}

// registerRoutes sets up the HTTP routes.
func (s *APIServer) registerRoutes() {
	s.router.HandleFunc("/", s.landingPage)
	s.router.HandleFunc("/import", s.importHandler)
	s.router.HandleFunc("/export", s.exportPage)
	s.router.HandleFunc("/help", s.helpPage)

	// Serve static assets (like styles.css) from the static directory.
	fs := http.FileServer(http.Dir("static"))
	s.router.Handle("/static/", http.StripPrefix("/static/", fs))
}

// landingPage renders the landing page.
func (s *APIServer) landingPage(w http.ResponseWriter, r *http.Request) {
	if err := s.templates.ExecuteTemplate(w, "landing.html", nil); err != nil {
		http.Error(w, "Error rendering landing page", http.StatusInternalServerError)
	}
}

// importHandler serves the import page on GET and processes form submissions on POST.
func (s *APIServer) importHandler(w http.ResponseWriter, r *http.Request) {
	// Debug log: record that the handler was called.
	log.Printf("importHandler called - Method: %s, URL: %s", r.Method, r.URL.String())

	if r.Method == http.MethodGet {
		if err := s.templates.ExecuteTemplate(w, "import.html", nil); err != nil {
			http.Error(w, "Error rendering import page", http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		log.Printf("POST request received at /import")

		// Parse the multipart form (up to 10MB)
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Error parsing form data", http.StatusInternalServerError)
			log.Printf("Error parsing form data: %v", err)
			return
		}

		// Check for file upload.
		file, handler, err := r.FormFile("datafile")
		if err == nil {
			defer file.Close()
			log.Printf("Uploaded file: %s", handler.Filename)
			// Process the file as needed.
		} else {
			log.Printf("No file uploaded or error: %v", err)
		}

		// Retrieve pasted data.
		pastedData := r.FormValue("data")
		log.Printf("Pasted data: %s", pastedData)

		// Process and sort the data as needed.
		// Render the success template.
		data := struct {
			Message string
		}{
			Message: "Data Imported Successfully!",
		}
		if err := s.templates.ExecuteTemplate(w, "success.html", data); err != nil {
			http.Error(w, "Error rendering success page", http.StatusInternalServerError)
			log.Printf("Error rendering success template: %v", err)
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Printf("Method not allowed: %s", r.Method)
	}
}

// exportPage renders the export page.
func (s *APIServer) exportPage(w http.ResponseWriter, r *http.Request) {
	if err := s.templates.ExecuteTemplate(w, "export.html", nil); err != nil {
		http.Error(w, "Error rendering export page", http.StatusInternalServerError)
	}
}

// helpPage renders the help page.
func (s *APIServer) helpPage(w http.ResponseWriter, r *http.Request) {
	if err := s.templates.ExecuteTemplate(w, "help.html", nil); err != nil {
		http.Error(w, "Error rendering help page", http.StatusInternalServerError)
	}
}

// Run starts the server.
func (s *APIServer) Run() {
	s.registerRoutes()
	log.Println("Server running on port:", s.serverPort)
	if err := http.ListenAndServe(s.serverPort, s.router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
