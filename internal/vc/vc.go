package vc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"time"
)

// Credential is a basic Verifiable Credential structure.
type Credential struct {
	Context           string                 `json:"@context"`
	Type              []string               `json:"type"`
	Issuer            string                 `json:"issuer"`
	IssuanceDate      string                 `json:"issuanceDate"`
	CredentialSubject map[string]interface{} `json:"credentialSubject"`
	Proof             string                 `json:"proof"`
}

// IssueDelegation creates and signs a simple credential asserting delegation.
func IssueDelegation(issuer, subjectDID string, metadata map[string]interface{}, secret []byte) (*Credential, error) {
	cred := &Credential{
		Context:      "https://www.w3.org/2018/credentials/v1",
		Type:         []string{"VerifiableCredential", "AgentDelegation"},
		Issuer:       issuer,
		IssuanceDate: time.Now().UTC().Format(time.RFC3339),
		CredentialSubject: map[string]interface{}{
			"id":       subjectDID,
			"metadata": metadata,
		},
	}

	payload, err := json.Marshal(cred)
	if err != nil {
		return nil, err
	}
	mac := hmac.New(sha256.New, secret)
	mac.Write(payload)
	cred.Proof = base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return cred, nil
}
