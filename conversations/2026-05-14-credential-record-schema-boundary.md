# Conversation: Credential Record Schema Boundary

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-define-credential-record-schema-boundary/`

Related artifacts:

- `docs/credential-record-schema-boundary.md`
- `docs/credential-record-schema-boundary.zh-CN.md`
- `decisions/ADR-0032-credential-record-schema-boundary.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`

## Context

The maintainer asked the agent to continue advancing bounded work items unless a true decision or blocker required confirmation. W-0074 follows the closeout of M-013 and starts M-014 credential/token schema ratification.

## Maintainer Narrative

The maintainer wants vibit to remain self-bootstrapping and controlled. Preparation should be thorough because the project is intentionally agent-native, but routine technical decisions inside the accepted direction should be handled by the agent's professional judgment instead of stopping for confirmation.

## Agent Response Summary

The agent ratified the credential record schema boundary for the selected `device_credential_login` posture.

The boundary defines future credential record ownership, lifecycle states, verifier semantics, uniqueness rules, player account relationship, rotation and replacement behavior, disabled/revoked credential behavior, redaction rules, and future migration/repository/adapter gates.

The agent did not add migration source, tables, repository interfaces, PostgreSQL adapters, runtime lookup, login handlers, token behavior, generated authentication shapes, Protobuf messages, WebSocket proof carriers, WebSocket handshake authentication, authentication dependencies, or production authentication behavior.

## Decisions

- Use `runtime.authentication` as the credential record owner.
- Ratify `authentication_device_credentials` as the future logical table name.
- Treat raw credential proof and raw operating-system device IDs as forbidden storage.
- Require non-plaintext, versioned credential lookup and verifier material.
- Treat only `credential_record_id` as log-safe.
- Ratify `active`, `disabled`, `revoked`, and `replaced` lifecycle states.
- Keep `player_id` immutable on credential records.
- Forbid moving a credential between players.
- Allow at most one active device credential per player in the first posture.
- Allow historical credential records for rotation, replacement, and revocation auditability.
- Keep multi-device linking, account recovery, and account merge deferred.
- Make W-0075 token verifier record schema boundary the next ready work item.

## Artifacts

- `docs/credential-record-schema-boundary.md`
- `docs/credential-record-schema-boundary.zh-CN.md`
- `decisions/ADR-0032-credential-record-schema-boundary.md`
- `changes/2026-05-14-define-credential-record-schema-boundary/`

## Open Questions

- Exact token verifier record fields, verifier algorithm, credential-token linkage, retention, cleanup, and revocation semantics.
- Exact migration order after both credential and token verifier schema boundaries are ratified.
- Exact repository interface shape for future credential lookup and mutation.
- Exact verifier algorithm and secret configuration boundary for runtime implementation.

## Follow-Up

- Advance W-0075.
- Keep migrations, repositories, adapters, runtime lookup, handlers, routes, generated authentication output, Protobuf changes, WebSocket changes, authentication dependencies, and authentication implementation behind later gates.

## Redaction Notes

No secrets, tokens, credential values, account details, provider payloads, or private data are stored in this conversation log.
