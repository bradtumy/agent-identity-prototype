# Agent Identity Prototype

This project is an architecture and working prototype for secure agent identity, delegation, and credential issuance using modern decentralized identity standards.

## Core Capabilities
- Agent registration via authenticated user
- Agent identity issuance using Verifiable Credentials (VC)
- Delegation enforcement via signed agent tokens
- Integration with OIDC + DID methods

## Tech Stack
- Go (Broker, Agent services)
- Keycloak (OIDC Identity Provider)
- Docker / Docker Compose
- DIDKit (for issuing DIDs and VCs)
- Optionally: Reused modules from [`go-oidc4vc-demo`](https://github.com/tumy-tech-labs/go-oidc4vc-demo)

## Architecture
![Architecture](proto-docs/architecture-diagram.png)

## How to Run
```bash
make init
make docker-up
```

Then:

Access Keycloak: http://localhost:8080
Run broker: go run broker/main.go
Use Postman to authenticate and register agents


