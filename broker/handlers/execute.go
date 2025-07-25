package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bradtumy/agent-identity-poc/internal/audit"
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
func ExecuteHandler(signingSecret []byte) http.HandlerFunc {
	// trusted issuer list and shared secret used for VC validation
	trustedIssuers := []string{"http://keycloak:8080/realms/agent-identity-poc"}
	sharedSecret := []byte("mysecret")
	return func(w http.ResponseWriter, r *http.Request) {
		var req ExecuteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("failed to decode execute request:", err)
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}

		cred := req.Credential

		if err := vc.VerifySignature(&cred, sharedSecret); err != nil {
			subj := cred.CredentialSubject.ID
			audit.LogAction("execute", subj, false)
			http.Error(w, "invalid credential signature", http.StatusUnauthorized)
			return
		}

		if err := vc.CheckTrustedIssuer(&cred, trustedIssuers); err != nil {
			subj := cred.CredentialSubject.ID
			audit.LogAction("execute", subj, false)
			http.Error(w, "untrusted issuer", http.StatusUnauthorized)
			return
		}

		if err := vc.CheckTTL(&cred); err != nil {
			subj := cred.CredentialSubject.ID
			audit.LogAction("execute", subj, false)
			http.Error(w, "expired credential", http.StatusUnauthorized)
			return
		}

		meta := cred.CredentialSubject.Metadata
		role, roleOK := meta["role"].(string)
		if !roleOK {
			subj := cred.CredentialSubject.ID
			audit.LogAction("execute", subj, false)
			http.Error(w, "missing role", http.StatusForbidden)
			return
		}

		action := req.Task.Action
		if err := policy.ValidatePolicy(action, role); err != nil {
			subj := cred.CredentialSubject.ID
			audit.LogAction("execute", subj, false)
			http.Error(w, "policy check failed: "+err.Error(), http.StatusForbidden)
			return
		}

		// Log success
		subj := cred.CredentialSubject.ID
		audit.LogAction("execute", subj, true)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"result": "ok"})
	}
}
