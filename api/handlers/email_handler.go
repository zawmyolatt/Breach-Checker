package handlers

import (
	"breach-checker/api/cache"
	"breach-checker/api/database"
	"breach-checker/api/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"time"
)

// EmailHandler handles email-related requests
type EmailHandler struct {
	db    *sql.DB
	cache *cache.RedisClient
}

// NewEmailHandler creates a new email handler
func NewEmailHandler(db *sql.DB, cache *cache.RedisClient) *EmailHandler {
	return &EmailHandler{
		db:    db,
		cache: cache,
	}
}

// CheckEmail checks if an email has been compromised
func (h *EmailHandler) CheckEmail(w http.ResponseWriter, r *http.Request) {
	// Set content type
	w.Header().Set("Content-Type", "application/json")

	var email string

	// Handle different HTTP methods
	switch r.Method {
	case "POST":
		// Parse request body for POST
		var req models.EmailCheckRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		email = strings.TrimSpace(req.Email)
	case "GET":
		// Get email from query parameter for GET
		email = strings.TrimSpace(r.URL.Query().Get("email"))
	default:
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Validate email
	if email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required")
		return
	}

	// Validate email format
	if _, err := mail.ParseAddress(email); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid email format")
		return
	}

	// Check cache first
	ctx := r.Context()
	cacheKey := "email:" + email
	cachedResult, err := h.cache.Get(ctx, cacheKey)

	if err == nil {
		// Cache hit
		var response models.EmailCheckResponse
		if err := json.Unmarshal([]byte(cachedResult), &response); err == nil {
			log.Printf("Cache hit for email: %s", email)
			respondWithJSON(w, http.StatusOK, response)
			return
		}
	}

	// Check database
	compromised, err := database.CheckEmailCompromised(h.db, email)
	if err != nil {
		log.Printf("Error checking email: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error checking email")
		return
	}

	// Prepare response
	response := models.EmailCheckResponse{
		Email:       email,
		Compromised: compromised,
	}

	if compromised {
		response.Message = "This email appears in our database of compromised accounts. We recommend changing your password immediately."
	} else {
		response.Message = "This email does not appear in our database of compromised accounts."
	}

	// Cache the result for 1 hour
	responseJSON, _ := json.Marshal(response)
	h.cache.Set(ctx, cacheKey, responseJSON, time.Hour)

	// Send response
	respondWithJSON(w, http.StatusOK, response)
}

// respondWithError returns an error response
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithJSON returns a JSON response
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Add a health check handler
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}
