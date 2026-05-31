# vibit Examples

Status: Draft v0.1
Last updated: 2026-05-31

The paired Simplified Chinese translation is `examples/README.zh-CN.md`. The English file is authoritative.

This directory contains source-first local examples and templates. These files are not product SDKs, release artifacts, hosted demos, install scripts, package registry publications, or direct Nakama/Pitaya compatibility surfaces.

## Local Alpha Example Client Path

Run from the repository root:

```bash
examples/local-alpha-example-client.sh
```

The source-first guide is:

```text
examples/local-alpha-client/README.md
```

The script wraps the focused authenticated gameplay E2E proof:

```bash
cd runtime && go test ./internal/platform/protocol/protobuf -run 'TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout|TestStorageObjectsProtocolRouteLocalAlphaFlow|TestFriendsRelationshipProtocolRouteLocalAlphaFlow|TestPresenceStatusLocalAlphaFlowReportsOfflineAfterCloseAndInvalidation|TestAuthenticatedGameplayFailurePathsLocalAlphaFlow' -v
```

It proves:

```text
local onboarding
-> device credential login
-> first-message connection binding
-> protected inventory grant/read
-> protected presence query
-> protected own-player storage object put/get/list/delete
-> protected friends send/status/accept/list/remove/block/unblock/reject
-> presence online/offline proof after close and invalidation
-> logout
-> post-logout protected request rejection
-> protected request failure-path and redaction checks
```

The storage object proof uses the existing `storage.GetOwnStorageObject`, `storage.ListOwnStorageObjects`, `storage.PutOwnStorageObject`, and `storage.DeleteOwnStorageObject` routes with `vibit.storage.v1` Protobuf payloads. It demonstrates Nakama-class durable player storage object capability coverage at the local request-flow level while preserving vibit's own route names and no direct Nakama/Pitaya API compatibility. It also demonstrates Pitaya-aligned layering by keeping transport, protocol adaptation, session metadata, route protection, application handlers, service behavior, and repository handoff separate.

The friends relationship proof uses the existing `friends.SendFriendRequest`, `friends.AcceptFriendRequest`, `friends.RejectFriendRequest`, `friends.RemoveFriend`, `friends.BlockPlayer`, `friends.UnblockPlayer`, `friends.ListFriendRelationships`, and `friends.GetFriendRelationshipStatus` routes with `vibit.friends.v1` Protobuf payloads. It demonstrates Nakama-class player social graph capability coverage at the local request-flow level while preserving vibit's own route names and no direct Nakama/Pitaya API compatibility. It also demonstrates Pitaya-aligned layering by keeping transport, protocol adaptation, session metadata, route protection, application handlers, service behavior, and repository handoff separate.

The script is intentionally redacted. It must not print raw credentials, raw access tokens, verifier keys, DSNs, digests, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or concrete transport metadata.

The older `examples/local-alpha-request-loop.sh` entrypoint remains as a compatibility wrapper around `examples/local-alpha-example-client.sh`.

## Local Environment Template

`examples/local.prototype.env.example` is a placeholder checklist for local configuration fields. It contains no real secrets and is not meant to be sourced directly without replacing every placeholder with local-only values.

Private local env files such as `.vibit.local.env`, `.env.local`, and `.env.*.local` are ignored by the repository `.gitignore`. Do not commit private local configuration.
