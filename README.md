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

1. **Authenticate to Keycloak** to obtain an access token. The default realm
   contains a public client `agent-identity-cli` and two demo users. Retrieve a
   token using the password grant:

   ```bash
   curl -X POST \
     -d 'client_id=agent-identity-cli' \
     -d 'grant_type=password' \
     -d 'username=alice' \
     -d 'password=password' \
     http://localhost:8080/realms/agent-identity-poc/protocol/openid-connect/token
   ```

   The JSON response contains an `access_token` field.

2. **Call the broker** with the obtained token:

   ```bash
   curl -X POST http://localhost:8081/register-agent \
     -H "Authorization: Bearer <access_token>" \
     -H "Content-Type: application/json" \
     -d '{"role":"data-fetcher","token_ttl":3600}'
   ```

   On success the broker returns the generated DID and a signed delegation
   credential. An example response is shown below:

   ```json
   {
     "did": "did:example:123",
     "credential": {
       "@context": "https://www.w3.org/2018/credentials/v1",
       "type": [
         "VerifiableCredential",
         "AgentDelegation"
       ],
       "issuer": "http://localhost:8081",
       "issuanceDate": "2025-07-24T13:52:35Z",
       "credentialSubject": {
         "id": "did:example:123",
         "metadata": {
           "role": "data-fetcher",
           "token_ttl": 3600
         }
       },
       "proof": "/rL2t/Ch9aklOZaV5fuamV/RwEfiuO/EfW5rBlNiL6k="
     }
   }
   ```
