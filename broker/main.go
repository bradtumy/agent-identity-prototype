package main

import (
	"log"
	"net/http"
	"os"

	"agent-broker/handlers"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Register routes
	r.HandleFunc("/register-agent", handlers.RegisterAgentHandler).Methods("POST")

	// Start server
	port := os.Getenv("BROKER_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Delegation Broker running on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
