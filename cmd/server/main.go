package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/deepto98/cargo-fullstack/internal/api"
	"github.com/deepto98/cargo-fullstack/internal/db"
	"github.com/deepto98/cargo-fullstack/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found: %v", err)
		// Continue execution as env vars might be set in the system
	}
	// Read DB connection string from environment or flag.
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// For local testing, adjust as necessary.
		// connStr = "postgresql://neondb_owner:npg_nuOB4RyQtCJ2@ep-raspy-resonance-a82l7j9d-pooler.eastus2.azure.neon.tech/polling_app?sslmode=require"
	}

	// Initialize the database.
	database, err := db.InitDB(connStr)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer database.Close()

	// Create API handler.
	apiHandler := &api.Handler{DB: database}

	// Load the frontend template.
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	// Create router.
	router := mux.NewRouter()

	// Serve static assets (CSS, JS, etc.)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// API endpoints.
	router.HandleFunc("/api/polls", apiHandler.CreatePoll).Methods("POST")
	router.HandleFunc("/api/polls/{id:[0-9]+}", apiHandler.GetPoll).Methods("GET")
	router.HandleFunc("/api/polls/{id:[0-9]+}/vote", apiHandler.Vote).Methods("POST")

	// Frontend: Serve index page.
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Render the main page.
		tmpl.Execute(w, nil)
	}).Methods("GET")

	// Wrap with logging middleware.
	loggedRouter := middleware.LoggingMiddleware(router)

	// Poll page route.
	router.HandleFunc("/poll/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("templates/poll.html")
		if err != nil {
			http.Error(w, "Template parsing error", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	}).Methods("GET")

	// Set up server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Handler: loggedRouter,
		Addr:    ":" + port,
	}

	// Start server in a goroutine.
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown.
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}
	log.Println("Server stopped gracefully")
}
