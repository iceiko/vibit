# Checklist

## Requirement

- [x] Requirement clarified.
- [x] Nakama/Pitaya alignment recorded.
- [x] Non-goals recorded.
- [x] Acceptance criteria recorded.

## Architecture

- [x] Affected modules identified.
- [x] Ownership impact reviewed.
- [x] Protocol boundary reviewed.
- [x] Generated-output standard reviewed.
- [x] Authentication/session boundary reviewed.
- [x] Repository, adapter, migration, dependency, blob/S3, hosted, release, announcement, promotion, ACL, and direct compatibility deferrals preserved.

## Implementation

- [x] Protobuf source added.
- [x] Generated Go Protobuf output regenerated.
- [x] Route keys added.
- [x] Protocol bridge mapping added.
- [x] Application route handlers added.
- [x] Startup registration added for PostgreSQL runtime path.
- [x] Transaction bypass added for service-owned write routes.

## Tests

- [x] Route handler tests added.
- [x] Protocol bridge tests added.
- [x] Protected route tests updated.
- [x] Startup id generation test added.
- [x] Focused Go package tests passed.

## Documentation

- [x] ADR added.
- [x] Conversation log added.
- [x] Change spec files added.
- [x] Manifests and AGENTS guides updated.
- [x] Check rules updated.

## Verification

- [x] Verification run or explicitly recorded.
- [x] Results recorded in `verification.md`.
- [ ] Commit and push.
