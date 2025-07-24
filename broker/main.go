package main

import (
	"context"
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
	storePath := getenv("STORAGE_PATH", "data/agents.json")

	store := storage.NewFileStore(storePath)

	auth, err := middleware.NewAuth(context.Background(), issuer)
	if err != nil {
		log.Fatalf("auth middleware init failed: %v", err)
	}

	r := mux.NewRouter()

	r.Handle("/register-agent", auth.Middleware(handlers.RegisterAgentHandler(store, issuer, signingSecret))).Methods("POST")

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
