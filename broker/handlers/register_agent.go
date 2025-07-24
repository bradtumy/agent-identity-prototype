package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bradtumy/agent-identity-poc/internal/did"
	"github.com/bradtumy/agent-identity-poc/internal/storage"
	"github.com/bradtumy/agent-identity-poc/internal/vc"
)

// AgentRequest represents the expected registration payload.
type AgentRequest struct {
	Role     string `json:"role"`
	TokenTTL int    `json:"token_ttl"`
}

// Response contains the issued credential.
type Response struct {
	DID        string         `json:"did"`
	Credential *vc.Credential `json:"credential"`
}

// RegisterAgentHandler handles POST /register-agent
func RegisterAgentHandler(store *storage.FileStore, issuer string, signingSecret []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// userEmail is set by auth middleware
		email, ok := r.Context().Value("userEmail").(string)
		if !ok || email == "" {
			http.Error(w, "missing user email", http.StatusUnauthorized)
			return
		}

		var req AgentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		agentDID := did.Generate()

		metadata := map[string]interface{}{
			"role":      req.Role,
			"token_ttl": req.TokenTTL,
		}

		cred, err := vc.IssueDelegation(issuer, agentDID, metadata, signingSecret)
		if err != nil {
			log.Printf("credential issuance error: %v", err)
			http.Error(w, "failed to issue credential", http.StatusInternalServerError)
			return
		}

		err = store.Save(storage.Agent{
			DID:        agentDID,
			Owner:      email,
			Metadata:   metadata,
			Credential: cred,
		})
		if err != nil {
			log.Printf("storage error: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(Response{DID: agentDID, Credential: cred})
	}
}
