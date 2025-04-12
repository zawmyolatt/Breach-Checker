package main

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

// EmailCheckRequest represents a request to check an email
type EmailCheckRequest struct {
	Email string `json:"email"`
}

// EmailCheckResponse represents a response to an email check
type EmailCheckResponse struct {
	Email       string `json:"email"`
	Compromised bool   `json:"compromised"`
	Message     string `json:"message,omitempty"`
}

func main() {
	// Create router
	r := mux.NewRouter()
	
	// Load templates
	templates := template.Must(template.ParseGlob("templates/*.html"))
	
	// Serve static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	
	// Home page
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		templates.ExecuteTemplate(w, "index.html", nil)
	}).Methods("GET")
	
	// Check email
	r.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		// Parse form
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		
		// Get email from form
		email := r.FormValue("email")
		if email == "" {
			http.Error(w, "Email is required", http.StatusBadRequest)
			return
		}
		
		// Create request to API
		apiURL := os.Getenv("API_URL")
		if apiURL == "" {
			apiURL = "http://api:8080"
		}
		
		reqBody, _ := json.Marshal(EmailCheckRequest{Email: email})
		req, err := http.NewRequest("POST", apiURL+"/api/check", bytes.NewBuffer(reqBody))
		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return
		}
		
		req.Header.Set("Content-Type", "application/json")
		
		// Send request to API
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Error connecting to API", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		
		// Parse response
		var result EmailCheckResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			http.Error(w, "Error parsing API response", http.StatusInternalServerError)
			return
		}
		
		// Render template with result
		templates.ExecuteTemplate(w, "result.html", result)
	}).Methods("POST")
	
	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")
	
	// Create server
	srv := &http.Server{
		Addr:         "0.0.0.0:3000",
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	// Start server in a goroutine
	go func() {
		log.Println("Web server starting on port 3000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	
	// Wait for interrupt signal to gracefully shut down the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	
	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	
	// Doesn't block if no connections, but will otherwise wait until the timeout
	srv.Shutdown(ctx)
	log.Println("Server gracefully stopped")
} 