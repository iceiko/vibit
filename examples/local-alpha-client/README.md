# Local Alpha Example Client Path

Status: Draft v0.1
Last updated: 2026-05-24

The paired Simplified Chinese translation is `README.zh-CN.md`. The English file is authoritative.

This directory records the first source-first repository-local local alpha example client path for vibit. It is a repository-local example path, not a public SDK, generated client library, package publication, hosted demo, release artifact, install script, live external client guarantee, or direct Nakama/Pitaya API compatibility surface.

## Run

From the repository root:

```bash
examples/local-alpha-example-client.sh
```

The script runs the focused local alpha Protobuf E2E proof inside the Go runtime:

```bash
cd runtime && go test ./internal/platform/protocol/protobuf -run 'TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout|TestStorageObjectsProtocolRouteLocalAlphaFlow|TestPresenceStatusLocalAlphaFlowReportsOfflineAfterCloseAndInvalidation|TestAuthenticatedGameplayFailurePathsLocalAlphaFlow' -v
```

## Demonstrated Flow

The example path demonstrates the current source-first local alpha loop:

```text
local onboarding
-> device credential login
-> first-message connection binding
-> protected inventory grant/read
-> protected presence query
-> protected own-player storage object put/get/list/delete
-> presence online/offline proof after close and invalidation
-> logout
-> rejected post-logout protected request
-> protected request failure-path and redaction checks
```

It uses existing runtime and protocol behavior:

```text
runtime.authentication.AuthenticateWithDeviceCredential
runtime.connection.BindConnection
inventory.GrantItem
inventory.GetInventory
presence.GetPlayerPresence
storage.PutOwnStorageObject
storage.GetOwnStorageObject
storage.ListOwnStorageObjects
storage.DeleteOwnStorageObject
runtime.authentication.LogoutAccessToken
```

## What This Is

This path is intended for developers and AI agents who need to see the current alpha capability loop without reverse-engineering the internal E2E test first. It keeps the proof source-first because the current generated Protobuf Go output lives under `runtime/internal/`, and local onboarding is still an application-owned setup behavior rather than a public client route.

The path maps to Nakama's developer-experience pressure: a backend framework should make it easy to understand how a client exercises authentication, protected gameplay requests, player storage, presence, logout, and failure behavior. vibit adapts that pressure without copying Nakama public routes, payloads, SDK shapes, token semantics, runtime APIs, or compatibility promises.

## What This Is Not

This path does not:

- publish a client SDK;
- generate client libraries;
- add a public onboarding protocol route;
- add protocol routes or Protobuf messages;
- change generated output;
- change runtime, authentication, session, persistence, migration, startup, or transport behavior;
- add dependencies;
- add hosted deployment or release artifacts;
- add stream subscriptions, chat rooms, groups, broadcast fanout, delivery guarantees, matchmaking, match runtime, operations/admin behavior, or distributed runtime;
- add direct Nakama/Pitaya API compatibility.

Pitaya remains deferred as a future distributed architecture reference. This path does not introduce frontend/backend server roles, RPC, service discovery, groups, cluster routing, or distributed sessions.

## Redaction

The script and example docs must not print, persist, commit, or record:

- raw device credential text or bytes;
- raw access tokens;
- credential or token lookup digests;
- credential or token verifier digests;
- verifier key values;
- concrete verifier key set ids;
- HMAC inputs or outputs;
- PostgreSQL DSNs with credentials;
- headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or concrete transport metadata.

Allowed output is limited to route names, step names, test names, redacted status classes, and high-level success/failure descriptions.

## Verification

The example path itself is verified by:

```bash
examples/local-alpha-example-client.sh
```

Architecture and continuation verification should also include:

```bash
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.local_alpha_example_client_path_implementation
node tools/vibit check change implement-local-alpha-example-client-path --json
node tools/vibit check runtime --json
node tools/vibit check all --json
git diff --check
```
