package handlers

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

// DelegateRequest is the expected payload for delegation.
type DelegateRequest struct {
	DelegateeDID string `json:"delegatee_did"`
	Role         string `json:"role"`
	TokenTTL     int    `json:"token_ttl"`
}

// DelegationToken represents the signed delegation credential.
type DelegationToken struct {
	Issuer            string                 `json:"issuer"`
	CredentialSubject map[string]string      `json:"credentialSubject"`
	Metadata          map[string]interface{} `json:"metadata"`
	IssuanceDate      string                 `json:"issuanceDate"`
	Proof             string                 `json:"proof"`
}

// DelegateHandler handles POST /delegate requests.
func DelegateHandler(issuer string, privKey ed25519.PrivateKey) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DelegateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}
		if req.DelegateeDID == "" || req.Role == "" || req.TokenTTL <= 0 {
			http.Error(w, "missing fields", http.StatusBadRequest)
			return
		}

		token := DelegationToken{
			Issuer:            issuer,
			CredentialSubject: map[string]string{"id": req.DelegateeDID},
			Metadata: map[string]interface{}{
				"role":      req.Role,
				"token_ttl": req.TokenTTL,
			},
			IssuanceDate: time.Now().UTC().Format(time.RFC3339),
		}

		payload, err := json.Marshal(token)
		if err != nil {
			http.Error(w, "encoding error", http.StatusInternalServerError)
			return
		}

		sig := ed25519.Sign(privKey, payload)
		token.Proof = base64.StdEncoding.EncodeToString(sig)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(token)
	}
}
