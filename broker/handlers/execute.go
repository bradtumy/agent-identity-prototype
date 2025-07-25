package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bradtumy/agent-identity-poc/internal/audit"
	"github.com/bradtumy/agent-identity-poc/internal/vc"
)

// ExecuteRequest payload for POST /execute
type ExecuteRequest struct {
	Credential string      `json:"credential"`
	Task       TaskRequest `json:"task"`
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
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		cred := req.Credential

		if err := vc.Verify(&cred, signingSecret); err != nil {
			subj, _ := cred.CredentialSubject["id"].(string)
			audit.LogAction("execute", subj, false)
      
		var cred vc.Credential
		if err := json.Unmarshal([]byte(req.Credential), &cred); err != nil {
			http.Error(w, "invalid credential", http.StatusBadRequest)
			return
		}

		if err := vc.Verify(&cred, signingSecret); err != nil {
			subj, _ := cred.CredentialSubject["id"].(string)
			//audit.LogAction("execute", subj, false)
			audit.LogAction("execute", cred.CredentialSubject["id"].(string), false)
			http.Error(w, "invalid credential", http.StatusForbidden)
			return
		}

		meta, ok := cred.CredentialSubject["metadata"].(map[string]interface{})
		if !ok {
			http.Error(w, "invalid credential metadata", http.StatusForbidden)
			return
		}
		role, roleOK := meta["role"].(string)
		ttlFloat, ttlOK := meta["token_ttl"].(float64)
		issued, err := time.Parse(time.RFC3339, cred.IssuanceDate)
		if err != nil {
			http.Error(w, "invalid issuance date", http.StatusForbidden)
			return
		}
      
		if !ttlOK || time.Now().After(issued.Add(time.Duration(ttlFloat)*time.Second)) {
			subj, _ := cred.CredentialSubject["id"].(string)
			audit.LogAction("execute", subj, false)
			http.Error(w, "credential expired", http.StatusForbidden)
			return
		}
		if !roleOK || role != "data-fetcher" {
			subj, _ := cred.CredentialSubject["id"].(string)
			audit.LogAction("execute", subj, false)

		if time.Now().After(issued.Add(time.Duration(ttl) * time.Second)) {
			audit.LogAction("execute", cred.CredentialSubject["id"].(string), false)
			http.Error(w, "credential expired", http.StatusForbidden)
			return
		}
		if role != "data-fetcher" {
			audit.LogAction("execute", cred.CredentialSubject["id"].(string), false)

			http.Error(w, "unauthorized role", http.StatusForbidden)
			return
		}

		// Log success
		subj, _ := cred.CredentialSubject["id"].(string)
		audit.LogAction("execute", subj, true)

//		audit.LogAction("execute", subj, true)
		audit.LogAction("execute", cred.CredentialSubject["id"].(string), true)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"result": "ok"})
	}
}
