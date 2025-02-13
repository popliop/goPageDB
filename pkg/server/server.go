package server

import (
	"log"
	"net/http"

	"github.com/marsorm/goPageDB/pkg/db"
	"github.com/marsorm/goPageDB/pkg/tmpl"
	"github.com/marsorm/goPageDB/pkg/upload"
)

type TemplateMessage struct {
	Message string
}

func NewTemplateMessage(message string) TemplateMessage {
	return TemplateMessage{
		Message: message,
	}
}

type APIServer struct {
	port   string
	router *http.ServeMux
}

func NewAPIServer(port string) *APIServer {
	mux := http.NewServeMux()

	// Load templates from the static directory.
	if err := tmpl.LoadTemplates("static/*.html"); err != nil {
		log.Fatalf("Could not load templates: %v", err)
	}

	return &APIServer{
		port:   port,
		router: mux,
	}
}

func (s *APIServer) registerRoutes() {
	s.router.HandleFunc("/", s.landingPage)
	s.router.HandleFunc("GET /import", s.importHandlerGET)
	s.router.HandleFunc("POST /import", s.importHandlerPOST)
	s.router.HandleFunc("/export", s.exportHandler)
	s.router.HandleFunc("/help", s.helpPage)
	s.router.HandleFunc("/test", s.testHandler)

	// Serve static assets (e.g., CSS files) from the static directory.
	fs := http.FileServer(http.Dir("static"))
	s.router.Handle("/static/", http.StripPrefix("/static/", fs))
}

func (s *APIServer) landingPage(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.RenderTemplate(w, "landing", nil); err != nil {
		log.Printf("landingPage error: %v", err)
		http.Error(w, "Error rendering landing page", http.StatusInternalServerError)
	}
}

func (s *APIServer) importHandlerPOST(w http.ResponseWriter, r *http.Request) {
	log.Println("POST request received at /import")

	fileUpload, err := upload.HandleForm(r)
	if err != nil {
		log.Printf("Error handling file and data: %v", err)
		if renderErr := tmpl.RenderTemplate(w, "error", NewTemplateMessage(err.Error())); renderErr != nil {
			log.Printf("Error rendering error template: %v", renderErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	dbConn, err := db.NewPostgresStorage()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		if renderErr := tmpl.RenderTemplate(w, "error", NewTemplateMessage(err.Error())); renderErr != nil {
			log.Printf("Error rendering error template: %v", renderErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	defer dbConn.Close()

	err = db.UpsertCSVData(dbConn, fileUpload.Data)
	if err != nil {
		log.Printf("Error upserting CSV data: %v", err)
		if renderErr := tmpl.RenderTemplate(w, "error", NewTemplateMessage(err.Error())); renderErr != nil {
			log.Printf("Error rendering error template: %v", renderErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if err := tmpl.RenderTemplate(w, "success", NewTemplateMessage("File processed and data upserted successfully")); err != nil {
		log.Printf("Error rendering success template: %v", err)
		http.Error(w, "Error rendering success page", http.StatusInternalServerError)
	}
}

func (s *APIServer) importHandlerGET(w http.ResponseWriter, r *http.Request) {
	log.Printf("importHandler called - Method: %s, URL: %s", r.Method, r.URL.String())

	err := tmpl.RenderTemplate(w, "import.html", nil)
	if err != nil {
		log.Printf("Error rendering import page: %v", err)
		http.Error(w, "Error rendering import page", http.StatusInternalServerError)
	}
}

func (s *APIServer) exportHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("exportHandler called - Method: %s, URL: %s", r.Method, r.URL.String())
	// fix export
}

func (s *APIServer) helpPage(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.RenderTemplate(w, "help", nil); err != nil {
		log.Printf("helpPage error: %v", err)
		http.Error(w, "Error rendering help page", http.StatusInternalServerError)
	}
}

func (s *APIServer) testHandler(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.RenderTemplate(w, "index", "woo"); err != nil {
		log.Printf("testHandler error: %v", err)
		http.Error(w, "Error rendering test page", http.StatusInternalServerError)
	}
}

func (s *APIServer) Run() {
	s.registerRoutes()
	log.Printf("Server running on port: %s", s.port)
	if err := http.ListenAndServe(s.port, s.router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
