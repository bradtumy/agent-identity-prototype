package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bradtumy/agent-identity-poc/internal/audit"
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
func ExecuteHandler(signingSecret []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ExecuteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("failed to decode execute request:", err)
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		cred := req.Credential

		if err := vc.Verify(&cred, signingSecret); err != nil {
			subj := cred.CredentialSubject.ID
			audit.LogAction("execute", subj, false)
			http.Error(w, "invalid credential", http.StatusForbidden)
			return
		}

		meta := cred.CredentialSubject.Metadata
		role, roleOK := meta["role"].(string)
		ttlFloat, ttlOK := meta["token_ttl"].(float64)
		issued, err := time.Parse(time.RFC3339, cred.IssuanceDate)
		if err != nil {
			http.Error(w, "invalid issuance date", http.StatusForbidden)
			return
		}
		if !ttlOK || time.Now().After(issued.Add(time.Duration(ttlFloat)*time.Second)) {
			subj := cred.CredentialSubject.ID
			audit.LogAction("execute", subj, false)
			http.Error(w, "credential expired", http.StatusForbidden)
			return
		}
		if !roleOK || role != "data-fetcher" {
			subj := cred.CredentialSubject.ID
			audit.LogAction("execute", subj, false)
			http.Error(w, "unauthorized role", http.StatusForbidden)
			return
		}

		// Log success
		subj := cred.CredentialSubject.ID
		audit.LogAction("execute", subj, true)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"result": "ok"})
	}
}
