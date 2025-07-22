package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// AgentIdentity is a mock return type
type AgentIdentity struct {
	ID       string `json:"id"`
	Owner    string `json:"owner"`
	Role     string `json:"role"`
	TokenTTL int    `json:"token_ttl"`
}

// RegisterAgentHandler handles POST /register-agent
func RegisterAgentHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}

	accessToken := strings.TrimPrefix(authHeader, "Bearer ")

	// 🔒 In real impl: validate JWT using Keycloak's JWKS
	log.Printf("Received token: %s\n", accessToken[:15]+"...")

	// Mock agent identity return
	identity := AgentIdentity{
		ID:       "agent-xyz123",
		Owner:    "user@example.com",
		Role:     "data-fetcher",
		TokenTTL: 3600,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(identity)
}
