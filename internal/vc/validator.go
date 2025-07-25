package vc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// VerifySignature validates the `proof` field using shared secret for now
func VerifySignature(cred *Credential, sharedSecret []byte) error {
	copyCred := *cred
	proof := copyCred.Proof
	copyCred.Proof = ""
	payload, err := json.Marshal(copyCred)
	if err != nil {
		return err
	}
	mac := hmac.New(sha256.New, sharedSecret)
	mac.Write(payload)
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(proof)) {
		return fmt.Errorf("invalid credential signature")
	}
	return nil
}

// CheckTrustedIssuer ensures the issuer matches known/trusted sources
func CheckTrustedIssuer(cred *Credential, trustedIssuers []string) error {
	for _, issuer := range trustedIssuers {
		if cred.Issuer == issuer {
			return nil
		}
	}
	return fmt.Errorf("untrusted issuer")
}

// CheckTTL ensures the credential has not expired
func CheckTTL(cred *Credential) error {
	issued, err := time.Parse(time.RFC3339, cred.IssuanceDate)
	if err != nil {
		return err
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
	default:
		return fmt.Errorf("invalid token_ttl type")
	}
	if time.Now().After(issued.Add(time.Duration(ttlSeconds) * time.Second)) {
		return fmt.Errorf("credential expired")
	}
	return nil
}
