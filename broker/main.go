package main

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bradtumy/agent-identity-poc/broker/handlers"
	"github.com/bradtumy/agent-identity-poc/broker/middleware"
	"github.com/bradtumy/agent-identity-poc/internal/storage"
	"github.com/gorilla/mux"
)

func main() {
	issuer := getenv("OIDC_ISSUER", "http://keycloak:8080/realms/agent-identity-poc")
	clientID := getenv("OIDC_CLIENT_ID", "agent-identity-cli")
	signingSecret := []byte(getenv("BROKER_SIGNING_SECRET", "secret"))
	keyB64 := getenv("BROKER_ED25519_PRIVATE_KEY", "")
	var privKey ed25519.PrivateKey
	if keyB64 != "" {
		keyBytes, err := base64.StdEncoding.DecodeString(keyB64)
		if err != nil {
			log.Fatalf("invalid ed25519 key: %v", err)
		}
		privKey = ed25519.PrivateKey(keyBytes)
	} else {
		_, pk, err := ed25519.GenerateKey(nil)
		if err != nil {
			log.Fatalf("key generation failed: %v", err)
		}
		privKey = pk
	}
	storePath := getenv("STORAGE_PATH", "data/agents.json")
	port := getenv("BROKER_PORT", "8081")

	log.Printf("Checking if OIDC issuer %s is ready...", issuer)
	if err := waitForOIDCIssuer(issuer, 10); err != nil {
		log.Fatalf("OIDC issuer not available: %v", err)
	}

	store := storage.NewFileStore(storePath)

	auth, err := middleware.NewAuth(context.Background(), issuer, clientID)
	if err != nil {
		log.Fatalf("auth middleware init failed: %v", err)
	}

	r := mux.NewRouter()
  
	r.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	r.Handle("/register-agent", auth.Middleware(handlers.RegisterAgentHandler(store, issuer, signingSecret))).Methods(http.MethodPost)
	r.Handle("/delegate", auth.Middleware(handlers.DelegateHandler(issuer, privKey))).Methods(http.MethodPost)

	port := getenv("BROKER_PORT", "8081")
  
	log.Printf("Delegation Broker running on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// waitForOIDCIssuer polls the OIDC metadata endpoint until it’s ready.
func waitForOIDCIssuer(issuer string, retries int) error {
	url := issuer + "/.well-known/openid-configuration"
	for i := 1; i <= retries; i++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			log.Println("OIDC issuer is ready.")
			return nil
		}
		if err != nil {
			log.Printf("OIDC check error: %v", err)
		} else {
			log.Printf("OIDC metadata returned status: %d", resp.StatusCode)
		}
		log.Printf("Waiting for OIDC issuer (%d/%d)...", i, retries)
		time.Sleep(5 * time.Second)
	}
	return fmt.Errorf("OIDC issuer %s not reachable after %d attempts", issuer, retries)
}

func getenv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
