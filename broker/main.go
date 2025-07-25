package main

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"log"
	"net/http"
	"os"

	"github.com/bradtumy/agent-identity-poc/broker/handlers"
	"github.com/bradtumy/agent-identity-poc/broker/middleware"
	"github.com/bradtumy/agent-identity-poc/internal/storage"
	"github.com/gorilla/mux"
)

func main() {
	issuer := getenv("OIDC_ISSUER", "http://localhost:8080/realms/agent-identity-poc")
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

	store := storage.NewFileStore(storePath)

	auth, err := middleware.NewAuth(context.Background(), issuer)
	if err != nil {
		log.Fatalf("auth middleware init failed: %v", err)
	}

	r := mux.NewRouter()

	r.Handle("/register-agent", auth.Middleware(handlers.RegisterAgentHandler(store, issuer, signingSecret))).Methods("POST")
	r.Handle("/delegate", auth.Middleware(handlers.DelegateHandler(issuer, privKey))).Methods("POST")

	port := getenv("BROKER_PORT", "8081")
	log.Printf("Delegation Broker running on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func getenv(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
