package models

import (
	"time"
)

// CompromisedEmail represents a compromised email in the database
type CompromisedEmail struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	BreachDate   time.Time `json:"breach_date"`
	BreachSource string    `json:"breach_source,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

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