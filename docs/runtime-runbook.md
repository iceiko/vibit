# Runtime Runbook

Status: Draft v0.2
Last updated: 2026-05-21
Scope: First Go runtime process startup, local alpha path, and manual verification

This runbook records how to start the first vibit Go runtime process and how the current local authenticated gameplay alpha path is shaped.

The paired Simplified Chinese translation is `docs/runtime-runbook.zh-CN.md`. The English file is authoritative.

## Current Runtime Surface

The current runtime process mounts one gameplay WebSocket endpoint:

```text
/v1/ws
```

The endpoint expects binary WebSocket messages containing `vibit.protocol.v1.Envelope` Protobuf bytes.

Text WebSocket messages are rejected by the transport adapter. JSON is not accepted on this endpoint.

The process also mounts small JSON troubleshooting endpoints:

```text
/healthz
/readyz
/version
/configz
```

`/healthz` reports whether the process is alive. `/readyz` reports the ready status, selected runtime store posture, and WebSocket path. `/version` reports the pre-alpha runtime version. `/configz` reports only redacted posture: runtime store, WebSocket path, local alpha request-loop script path, PostgreSQL configuration presence, whether authentication configuration is required for the selected composition, and `secrets_redacted: true`.

These endpoints are local alpha troubleshooting surfaces, not a production operations API, admin console, metrics backend, protocol route, or release packaging surface.

There are two startup compositions:

```text
VIBIT_RUNTIME_STORE=memory
```

The memory store is the default bootstrap path. It wires the original in-memory inventory request loop only:

```text
WebSocket binary frame
-> Protobuf envelope
-> application dispatch
-> inventory command or query handler
-> Protobuf response envelope
-> WebSocket binary frame
```

It does not wire the current authentication service, route protection, connection binding, logout, runtime sessions, or presence query.

```text
VIBIT_RUNTIME_STORE=postgres
```

The PostgreSQL store is the current alpha runtime composition. It wires:

- PostgreSQL-backed inventory state.
- Player account lifecycle persistence.
- Device credential login and opaque access-token issuance.
- Runtime session creation during login.
- Request-level access-token route protection through `AuthenticatedRequest`.
- First-message connection binding through `runtime.authentication.BindConnection`.
- Logout through `runtime.authentication.LogoutAccessToken`.
- Registry-backed presence lifecycle and the protected `runtime.presence.GetPlayerPresence` query.

Local onboarding now exists as an application service method, `OnboardLocalPlayerWithDeviceCredential`, but it is not a public WebSocket, Protobuf, HTTP, CLI, or startup auto-creation surface. It is proven by tests and reserved for future local tooling.

## Start The Server

From the Go runtime module:

```bash
cd runtime
go run ./cmd/vibit-server
```

The default listen address is:

```text
:8080
```

Override it with:

```bash
VIBIT_ADDR=:9090 go run ./cmd/vibit-server
```

The default runtime store is in memory:

```text
VIBIT_RUNTIME_STORE=memory
```

To start the explicit PostgreSQL-backed inventory composition path, provide both the store selector and a PostgreSQL DSN:

```bash
VIBIT_RUNTIME_STORE=postgres VIBIT_POSTGRES_DSN='postgres://user:pass@127.0.0.1:5432/vibit?sslmode=disable' go run ./cmd/vibit-server
```

The PostgreSQL runtime path also wires the current authentication, token, runtime session, route protection, connection binding, logout, and presence-lifecycle composition. It requires authentication verifier key environment variables:

```text
VIBIT_AUTH_VERIFIER_KEY_SET_ID
VIBIT_AUTH_CREDENTIAL_LOOKUP_KEY
VIBIT_AUTH_CREDENTIAL_VERIFIER_KEY
VIBIT_AUTH_TOKEN_LOOKUP_KEY
VIBIT_AUTH_TOKEN_VERIFIER_KEY
```

Verifier key values must be Base64 text accepted by the runtime loader. Do not commit local verifier keys.

Verifier key requirements:

- `VIBIT_AUTH_VERIFIER_KEY_SET_ID` must be non-empty.
- Each logical verifier key must decode to at least 32 bytes.
- The four logical verifier keys must be distinct.
- Weak repeated-byte keys are rejected.
- Key values and concrete key set ids are not log-safe. Treat them as secrets.

The runtime loader accepts URL-safe unpadded Base64 and standard padded Base64 key text. Keep verifier keys in local environment configuration or an explicit local secret source; do not commit them to the repository, shell history, ADRs, change records, test fixtures, or runbook examples.

Optional authentication settings:

```text
VIBIT_AUTH_ACCESS_TOKEN_TTL
VIBIT_AUTH_TOKEN_AUDIENCE
```

Optional PostgreSQL pool settings:

```text
VIBIT_POSTGRES_MAX_CONNS
VIBIT_POSTGRES_MIN_CONNS
```

Normal server startup does not apply migrations. Apply or verify migrations explicitly before using the PostgreSQL store path against a fresh database.

## Local Alpha Flow

The now-proven local alpha path is:

```text
local onboarding
-> device credential login
-> connection binding
-> protected inventory grant/read
-> protected presence query
-> logout
-> post-logout protected request rejection
```

The executable proof is:

```bash
cd runtime && go test ./internal/platform/protocol/protobuf -run TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout -v
```

The minimal local alpha request-loop script is:

```bash
examples/local-alpha-request-loop.sh
```

The script is a redacted wrapper around the focused E2E proof. It is not a public onboarding client, product SDK, live PostgreSQL process client, or release artifact.

That test uses the same protocol frame handler shape as the runtime for login, binding, protected inventory, protected presence, and logout. It uses test-local repositories, deterministic entropy, deterministic ids, and an in-memory registry so the flow can be verified without a live PostgreSQL server or committed secrets.

Current alpha path details:

1. Local onboarding creates an active player account and active digest-only device credential record through `OnboardLocalPlayerWithDeviceCredential`.
2. Local onboarding returns raw device credential text only once after unit-of-work success. It does not return an access token or runtime session.
3. `runtime.authentication.AuthenticateWithDeviceCredential` accepts the device credential proof and returns an opaque access token after validating the stored credential verifier.
4. The login response envelope carries runtime session metadata in `Envelope.Session`. The session id is metadata, not authentication proof.
5. `runtime.authentication.BindConnection` binds a server-observed connection id and epoch to the validated player identity using the access token. WebSocket handshake authentication is still not used.
6. Protected routes such as `inventory.GrantItem`, `inventory.GetInventory`, and `runtime.presence.GetPlayerPresence` require `vibit.authentication.v1.AuthenticatedRequest`.
7. `runtime.authentication.LogoutAccessToken` revokes only the presented access token.
8. A later protected request using the same revoked token is rejected.

The current public protocol does not include a local onboarding route. A developer who runs the PostgreSQL server process directly still needs a future local tool, request-loop script, or controlled seed path to call local onboarding and obtain the first device credential.

The packaged local alpha developer journey is:

```text
docs/alpha-developer-flow.md
```

It connects the README, this runbook, the redacted request-loop script, local status endpoints, the acceptance checklist, PostgreSQL manual setup posture, verification commands, and the next contribution path. The packaged flow exists now, `docs/release-publishing-decision-gate.md` defines the release publishing decision boundary, `docs/release-execution-preparation-gate.md` defines the release execution preparation boundary, `docs/release-execution-authorization-gate.md` defines the release execution authorization gate criteria, `docs/release-execution-maintainer-decision.md` records the maintainer go decision to continue planning, and `docs/release-identifier-artifact-plan.md` records the proposed `v0.1.0-alpha.1` identifier and source-first artifact plan. The next gap is final maintainer authorization, not release publishing or artifact creation itself.

The local alpha acceptance checklist is:

```text
docs/alpha-acceptance-checklist.md
```

It records which alpha items are ready, manual, deferred, or blocked. It is not a release declaration and does not authorize release packaging.

## Manual Verification Paths

### Bootstrap Memory Path

1. Start the server.
2. Connect a WebSocket client to `ws://127.0.0.1:8080/v1/ws`.
3. For the bootstrap in-memory path, send a binary Protobuf `Envelope` for `inventory.GrantItem` or `inventory.GetInventory`.
4. Confirm the response is a binary Protobuf `Envelope` with the same `request_id`.

### Local Alpha Request Loop

Run the minimal local request-loop script:

```bash
examples/local-alpha-request-loop.sh
```

The script prints only a redacted path summary and Go test status. It does not require live PostgreSQL, committed verifier keys, raw credentials, raw access tokens, or a hand-built WebSocket client.

### Authenticated Alpha Proof

Run the focused E2E test:

```bash
cd runtime && go test ./internal/platform/protocol/protobuf -run TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout -v
```

This is the current best proof of the complete local alpha flow. It does not require live PostgreSQL and does not print raw credentials or tokens.

### Alpha Acceptance Checklist

Read the current checklist:

```bash
sed -n '1,220p' docs/alpha-acceptance-checklist.md
```

Confirm the ready, manual, deferred, and blocked states before packaging or publishing work. The checklist keeps release publishing, release packaging, public local onboarding, production signup, broad operations/admin behavior, broad product modules, and direct Nakama/Pitaya API compatibility deferred.

### PostgreSQL Runtime Path

For the PostgreSQL process path:

1. Prepare a local PostgreSQL database.
2. Apply the SQL migrations explicitly. Normal server startup does not apply them.
3. Provide `VIBIT_RUNTIME_STORE=postgres`, `VIBIT_POSTGRES_DSN`, and all verifier key environment variables.
4. Start `go run ./cmd/vibit-server`.
5. Use binary Protobuf envelopes over `ws://127.0.0.1:8080/v1/ws`.

The process path currently exposes login, binding, protected inventory, protected presence, and logout routes. It does not expose local onboarding as a public protocol route.

## Current Runtime Assumptions

- The runtime uses an in-memory inventory repository by default.
- `VIBIT_RUNTIME_STORE=postgres` enables the explicit PostgreSQL composition path for persistent inventory, player account, authentication token/credential, runtime session, route protection, logout, connection binding, and presence-lifecycle wiring.
- Inventory bootstrap permissions allow grant and read operations.
- Local onboarding/device credential issuance exists as an application service method and is proven in tests. It is not public protocol behavior.
- The authenticated gameplay E2E path is proven by a focused Go protocol test over existing capabilities.
- `examples/local-alpha-request-loop.sh` is the minimal local alpha request-loop script over that proof.
- `/healthz`, `/readyz`, `/version`, and `/configz` provide the minimal local status surface for startup troubleshooting.
- PostgreSQL persistence is selected explicitly. The persistence boundary is defined in `docs/postgresql-persistence-boundary.md`.
- PostgreSQL migrations are not applied automatically during normal server startup.
- Optional live PostgreSQL verification is defined in `docs/postgresql-verification-environment.md`; it requires `VIBIT_POSTGRES_TEST_DSN` and is not part of default server startup.
- Generated route registration is not implemented yet; route registration remains handwritten startup/bootstrap code.
- The v0.1 alpha path has an alpha acceptance checklist at `docs/alpha-acceptance-checklist.md`.
- The v0.1 alpha path still needs the alpha developer flow packaged into one coherent local developer journey.
- Production signup, external identity providers, password login, account recovery, account merge, multi-device linking, direct Nakama/Pitaya API compatibility, and broad product modules remain deferred.

These are bootstrap assumptions for the first request loop, not long-term production policy.

## Redaction Rules

Do not record or commit:

- Raw device credential text or bytes.
- Raw access tokens.
- Credential or token lookup digests.
- Credential or token verifier digests.
- HMAC input or output bytes.
- Verifier key values.
- Concrete verifier key set ids.
- PostgreSQL DSNs with credentials.
- Headers, cookies, query strings, WebSocket subprotocol values, or remote addresses that may carry secrets.

Documentation, ADRs, change specs, logs, test output, and examples should use placeholders such as `redacted-token-text`, `redacted-device-credential`, or `postgres://user:pass@127.0.0.1:5432/vibit?sslmode=disable` only when the value is clearly non-production sample text.

## Verification Commands

Run from the repository root unless noted:

```bash
cd runtime && go test ./...
cd runtime && go test ./internal/platform/protocol/protobuf -run TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout -v
cd runtime && go vet ./...
node tools/vibit check runtime
node tools/vibit check postgres-env
node tools/vibit check all
```

`node tools/vibit check postgres-env` is a static standards check. It does not connect to PostgreSQL. Live PostgreSQL verification remains opt-in through `VIBIT_POSTGRES_TEST_DSN`.

Run the current live durable inventory verification against a disposable PostgreSQL database with:

```bash
cd runtime && VIBIT_POSTGRES_TEST_DSN='postgres://user:pass@127.0.0.1:5432/vibit_test?sslmode=disable' VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE=1 go test ./internal/platform/protocol/protobuf -run TestPostgresPersistentInventoryRequestLoop -v
```

This test applies the inventory migration explicitly and verifies the WebSocket Protobuf `GrantItem` then `GetInventory` request loop through the PostgreSQL-backed runtime composition. If `VIBIT_POSTGRES_TEST_DSN` is unset, the test skips and records that live PostgreSQL verification was unavailable.

The test uses `drop_schema` cleanup semantics by default. Other cleanup modes are intentionally skipped for this test because migration apply must be verified from a clean schema.
