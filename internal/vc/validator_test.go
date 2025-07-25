package vc

import (
	"testing"
	"time"
)

func TestVerifySignature(t *testing.T) {
	secret := []byte("mysecret")
	cred, err := IssueDelegation("http://keycloak:8080/realms/agent-identity-poc", "did:example:123", map[string]interface{}{"token_ttl": 3600}, secret)
	if err != nil {
		t.Fatalf("issue credential: %v", err)
	}
	if err := VerifySignature(cred, secret); err != nil {
		t.Fatalf("valid signature rejected: %v", err)
	}
	if err := VerifySignature(cred, []byte("wrong")); err == nil {
		t.Fatalf("invalid signature accepted")
	}
}

func TestCheckTrustedIssuer(t *testing.T) {
	secret := []byte("mysecret")
	issuer := "http://keycloak:8080/realms/agent-identity-poc"
	cred, _ := IssueDelegation(issuer, "did:example:123", map[string]interface{}{"token_ttl": 3600}, secret)
	if err := CheckTrustedIssuer(cred, []string{issuer}); err != nil {
		t.Fatalf("trusted issuer rejected: %v", err)
	}
	if err := CheckTrustedIssuer(cred, []string{"http://malicious"}); err == nil {
		t.Fatalf("untrusted issuer accepted")
	}
}

func TestCheckTTL(t *testing.T) {
	secret := []byte("mysecret")
	cred, _ := IssueDelegation("http://keycloak:8080/realms/agent-identity-poc", "did:example:123", map[string]interface{}{"token_ttl": 3600}, secret)
	if err := CheckTTL(cred); err != nil {
		t.Fatalf("valid ttl rejected: %v", err)
	}
	cred.IssuanceDate = time.Now().Add(-2 * time.Hour).Format(time.RFC3339)
	cred.CredentialSubject.Metadata["token_ttl"] = 1
	if err := CheckTTL(cred); err == nil {
		t.Fatalf("expired ttl accepted")
	}
}
