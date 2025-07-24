# Agent Identity Prototype

This project demonstrates secure agent identity creation and delegation using Go, Keycloak and simple DID/VC utilities.

## Core Capabilities
- Agent registration via authenticated user
- Agent identity issuance using Verifiable Credentials (VC)
- Delegation enforcement via signed agent tokens
- Integration with OIDC + DID methods

## Tech Stack
- Go (Broker, Agent services)
- Keycloak (OIDC Identity Provider)
- Docker / Docker Compose

## Running
```bash
make docker-up
```
This will start Keycloak, Redis, the delegation broker and a placeholder agent runner.
The broker listens on `http://localhost:8081`.

### Register an Agent
`POST /register-agent` with a Bearer token obtained from Keycloak and JSON body:
```json
{
  "role": "data-fetcher",
  "token_ttl": 3600
}
```
On success the broker returns the generated DID and a signed delegation credential.
