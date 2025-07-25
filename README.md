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

### Issue a Delegation Token

Use the `/delegate` endpoint to generate a signed delegation token for an agent.
This route **only accepts POST requests**.

```bash
curl -X POST http://localhost:8081/delegate \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "delegatee_did": "did:example:123",
    "role": "data-fetcher",
    "token_ttl": 3600
  }'
```

The response contains the signed token which includes the delegatee DID,
your role metadata and a proof signature.

The broker signs tokens using an Ed25519 private key. You may supply your own
key via the `BROKER_ED25519_PRIVATE_KEY` environment variable (base64 encoded).
If not provided, a new key is generated at startup.

### Execute a Task

The `/execute` endpoint allows an agent to perform an authorized action using
its delegation credential. The examples below show how to call it via
**Postman** and **curl**.

1. Create a new `POST` request to `http://localhost:8081/execute` in Postman.
2. On the **Body** tab choose **raw** and select **JSON** format.
3. Provide the payload containing the credential string and task details:

```json
{
  "credential": "<credential JSON string>",
  "task": {
    "action": "fetch_data",
    "params": {
      "url": "https://example.com/data"
    }
  }
}
```

Use the credential returned by `/register-agent`. Postman may require escaping quotes or using environment variables to include the full JSON string.

To call the endpoint with `curl`:

```bash
curl -X POST http://localhost:8081/execute \
  -H "Content-Type: application/json" \
  -d '{
    "credential": "<credential JSON string>",
    "task": {"action": "fetch_data", "params": {"url": "https://example.com/data"}}
  }'
```

Sending the request will return a stubbed result when the credential is valid:

```json
{"result": "ok"}
```

Each call writes a log entry to `./data/execution.log`. You should see a line
similar to:

```json
{"timestamp":"2025-07-25T15:42:00Z","agent_did":"did:example:abc-123","role":"data-fetcher","action":"fetch_data","status":"success","message":"Fetched data from https://example.com/data"}
```

Invalid or expired credentials receive a `403 Forbidden` response.


## Keycloak Configuration

When running `make docker-up` the Keycloak container automatically imports the
realm definition from `keycloak/realm-export.json`. This file enables *Direct
Access Grants* for the `agent-identity-cli` client and adds an audience mapper
so issued tokens contain the `aud` claim. If you start Keycloak manually,
import the same file using `--import-realm` or through the admin console.
You can verify this by requesting a token using the password grant:

```bash
curl -X POST \
  -d 'client_id=agent-identity-cli' \
  -d 'grant_type=password' \
  -d 'username=alice' \
  -d 'password=password' \
  http://localhost:8080/realms/agent-identity-poc/protocol/openid-connect/token
```

Decoding the `access_token` should show the audience claim injected by the protocol mapper:

```json
{
  "preferred_username": "alice",
  "aud": "agent-identity-cli"
}
```

This configuration is required so the broker and runner components can validate tokens issued to the CLI.

## Credential Trust Model

The `/execute` endpoint now enforces:

- ✅ Signature validation via `credential.proof` (HMAC for now)
- ✅ Trusted issuer check against known realm URL(s)
- ✅ TTL enforcement based on `issuanceDate` + `token_ttl`

If a credential fails any check, the server responds with 401 Unauthorized.

These settings will evolve to support DID + VC chains in future phases.


## Agent Execution Policy Engine

Before executing a task, the broker validates:

- Action is in the allowed list (`fetch_data`, `transform`, `notify`)
- Role in credential metadata is authorized for the requested action

If either check fails, a `403 Forbidden` is returned.

## Execution Logs

All agent task execution events are logged to `./data/execution.log` in JSON
Lines format.

