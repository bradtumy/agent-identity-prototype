package vc

import (
	"fmt"
	"time"
)

// ValidateTTL ensures the credential has not expired based on issuanceDate and token_ttl.
func ValidateTTL(cred *Credential) error {
	issued, err := time.Parse(time.RFC3339, cred.IssuanceDate)
	if err != nil {
		return fmt.Errorf("invalid issuanceDate: %w", err)
	}
	ttlVal, ok := cred.CredentialSubject.Metadata["token_ttl"]
	if !ok {
		return fmt.Errorf("missing token_ttl")
	}
	var ttlSeconds float64
	switch v := ttlVal.(type) {
	case float64:
		ttlSeconds = v
	case int:
		ttlSeconds = float64(v)
	case int64:
		ttlSeconds = float64(v)
	default:
		return fmt.Errorf("invalid token_ttl type")
	}
	expiry := issued.Add(time.Duration(ttlSeconds) * time.Second)
	if time.Now().UTC().After(expiry) {
		return fmt.Errorf("expired_token")
	}
	return nil
}
