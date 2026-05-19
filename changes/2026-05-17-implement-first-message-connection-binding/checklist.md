# Checklist

## Requirement

- [x] Requirement clarified.
- [x] Non-goals recorded.
- [x] Acceptance criteria recorded.

## Architecture

- [x] Affected modules identified.
- [x] Ownership impact reviewed.
- [x] Dependencies reviewed.
- [x] Nakama/Pitaya reference mapping preserved without direct API compatibility.

## Implementation

- [x] Protobuf source updated before generated output.
- [x] Generated Go Protobuf output updated through Buf.
- [x] Application binder implemented inside `runtime/internal/app`.
- [x] Protobuf adapter implemented inside `runtime/internal/platform/protocol/protobuf`.
- [x] Startup composition implemented inside `runtime/cmd/vibit-server`.
- [x] WebSocket transport remains credential-neutral.
- [x] Session persistence, migrations, repositories, logout, reconnect, and dependencies remain deferred.

## Tests

- [x] Application binder tests added.
- [x] Protocol adapter tests added.
- [x] Startup composition tests updated.
- [x] WebSocket transport tests updated.

## Documentation

- [x] Change specs added.
- [x] Conversation log added.
- [x] Work queue updated.
- [x] Manifests and AGENTS guides updated.
- [x] Repository check rule updated.

## Verification

- [x] Verification run.
- [x] Results recorded in `verification.md`.
