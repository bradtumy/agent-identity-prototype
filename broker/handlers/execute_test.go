package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bradtumy/agent-identity-poc/internal/vc"
)

func TestExecuteHandlerExpiredToken(t *testing.T) {
	secret := []byte("mysecret")
	cred, err := vc.IssueDelegation("http://keycloak:8080/realms/agent-identity-poc", "did:example:123", map[string]interface{}{"role": "data-fetcher", "token_ttl": 1}, secret)
	if err != nil {
		t.Fatalf("issue credential: %v", err)
	}

	time.Sleep(2 * time.Second)

	reqPayload := ExecuteRequest{Credential: *cred, Task: vc.Task{Action: "fetch_data"}}
	b, _ := json.Marshal(reqPayload)
	req := httptest.NewRequest(http.MethodPost, "/execute", bytes.NewReader(b))
	rec := httptest.NewRecorder()

	handler := ExecuteHandler(secret, nil)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", rec.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if resp["error"] != "expired_token" {
		t.Fatalf("unexpected response: %v", resp)
	}
}
