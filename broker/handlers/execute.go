package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bradtumy/agent-identity-poc/internal/audit"
	"github.com/bradtumy/agent-identity-poc/internal/executionlog"
	"github.com/bradtumy/agent-identity-poc/internal/policy"
	"github.com/bradtumy/agent-identity-poc/internal/vc"
)

// ExecuteRequest payload for POST /execute
type ExecuteRequest struct {
	Credential vc.Credential `json:"credential"`
	Task       vc.Task       `json:"task"`
}

// TaskRequest describes an agent action
type TaskRequest struct {
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"params"`
}

// ExecuteHandler handles POST /execute requests
func ExecuteHandler(signingSecret []byte, logger *executionlog.Logger) http.HandlerFunc {
	// trusted issuer list and shared secret used for VC validation
	trustedIssuers := []string{"http://keycloak:8080/realms/agent-identity-poc"}
	sharedSecret := []byte("mysecret")
	return func(w http.ResponseWriter, r *http.Request) {
		var req ExecuteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("failed to decode execute request:", err)
			http.Error(w, "invalid payload", http.StatusBadRequest)
			if logger != nil {
				entry := executionlog.Entry{
					Timestamp: time.Now().UTC().Format(time.RFC3339),
					Status:    "failure",
					Message:   "invalid payload",
				}
				if err := logger.Log(entry); err != nil {
					log.Printf("execution log error: %v", err)
				}
			}
			return
		}

		cred := req.Credential
		role, _ := cred.CredentialSubject.Metadata["role"].(string)
		action := req.Task.Action
		entry := executionlog.Entry{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			AgentDID:  cred.CredentialSubject.ID,
			Role:      role,
			Action:    action,
		}

		if err := vc.VerifySignature(&cred, sharedSecret); err != nil {
			subj := cred.CredentialSubject.ID
			audit.LogAction("execute", subj, false)
			entry.Status = "failure"
			entry.Message = "invalid credential signature"
			if logger != nil {
				if err := logger.Log(entry); err != nil {
					log.Printf("execution log error: %v", err)
				}
			}
			http.Error(w, "invalid credential signature", http.StatusUnauthorized)
			return
		}

		if err := vc.CheckTrustedIssuer(&cred, trustedIssuers); err != nil {
			subj := cred.CredentialSubject.ID
			audit.LogAction("execute", subj, false)
			entry.Status = "failure"
			entry.Message = "untrusted issuer"
			if logger != nil {
				if err := logger.Log(entry); err != nil {
					log.Printf("execution log error: %v", err)
				}
			}
			http.Error(w, "untrusted issuer", http.StatusUnauthorized)
			return
		}

		if err := vc.ValidateTTL(&cred); err != nil {
			subj := cred.CredentialSubject.ID
			log.Printf("token TTL validation failed for %s: %v", subj, err)
			audit.LogAction("execute", subj, false)
			entry.Status = "failure"
			entry.Message = "expired credential"
			if logger != nil {
				if err := logger.Log(entry); err != nil {
					log.Printf("execution log error: %v", err)
				}
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "expired_token",
				"message": "The delegation token has expired.",
			})
			return
		}

		meta := cred.CredentialSubject.Metadata
		role, roleOK := meta["role"].(string)
		if !roleOK {
			subj := cred.CredentialSubject.ID
			audit.LogAction("execute", subj, false)
			entry.Status = "failure"
			entry.Message = "missing role"
			if logger != nil {
				if err := logger.Log(entry); err != nil {
					log.Printf("execution log error: %v", err)
				}
			}
			http.Error(w, "missing role", http.StatusForbidden)
			return
		}

		if err := policy.ValidatePolicy(action, role); err != nil {
			subj := cred.CredentialSubject.ID
			audit.LogAction("execute", subj, false)
			entry.Status = "failure"
			entry.Message = "policy check failed: " + err.Error()
			if logger != nil {
				if err := logger.Log(entry); err != nil {
					log.Printf("execution log error: %v", err)
				}
			}
			http.Error(w, "policy check failed: "+err.Error(), http.StatusForbidden)
			return
		}

		// Log success
		subj := cred.CredentialSubject.ID
		audit.LogAction("execute", subj, true)
		entry.Status = "success"
		// Generate a simple success message
		successMsg := fmt.Sprintf("%s executed", action)
		if action == "fetch_data" {
			if url, ok := req.Task.Params["url"]; ok {
				successMsg = fmt.Sprintf("Fetched data from %v", url)
			}
		}
		entry.Message = successMsg
		if logger != nil {
			if err := logger.Log(entry); err != nil {
				log.Printf("execution log error: %v", err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"result": "ok"})
	}
}
