package vc

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
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
	return ValidateTTL(cred)
}
