# Conversation: Token Verifier Record Schema Boundary

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-define-token-verifier-record-schema-boundary/`

Related artifacts:

- `docs/token-verifier-record-schema-boundary.md`
- `docs/token-verifier-record-schema-boundary.zh-CN.md`
- `decisions/ADR-0033-token-verifier-record-schema-boundary.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`

## Context

The maintainer asked the agent to continue advancing bounded work items unless a true decision or blocker required confirmation. W-0075 follows W-0074 inside M-014 credential/token schema ratification.

## Maintainer Narrative

The maintainer wants vibit to remain self-bootstrapping and controlled. Routine technical decisions inside an already-ratified direction should be handled by the agent's professional judgment, while real product or architecture branch points should still be brought back for discussion.

## Agent Response Summary

The agent ratified the token verifier record schema boundary for the selected opaque access-token posture.

The boundary defines future token verifier record ownership, non-plaintext verifier storage, lifecycle states, expiration, revocation, logout scope, rotation linkage to credential records, retention, cleanup expectations, replay-sensitive failure classes, redaction rules, and future migration/repository/adapter gates.

The agent did not add migration source, tables, repository interfaces, PostgreSQL adapters, runtime token issuance, runtime token validation, logout, refresh, cleanup jobs, generated authentication shapes, Protobuf messages, WebSocket proof carriers, WebSocket handshake authentication, authentication dependencies, or production authentication behavior.

## Decisions

- Use `runtime.authentication` as the token verifier record owner.
- Ratify `authentication_access_tokens` as the future logical table name.
- Treat raw access-token storage as forbidden.
- Require non-plaintext, versioned token lookup and verifier material.
- Treat only `token_record_id` as the normal log-safe token identifier.
- Ratify `active`, `expired`, and `revoked` lifecycle states.
- Keep `player_id` immutable on token verifier records.
- Require `credential_record_id` linkage for the first posture.
- Forbid moving a token between players or credentials.
- Keep refresh-token storage, runtime session binding, WebSocket connection binding, Redis-like token/session stores, and service/admin token families deferred.
- Recommend seven-day retention for expired and revoked token verifier records until a later operations decision changes it.
- Require cleanup before production token storage, while deferring cleanup implementation.
- Make W-0076 authentication schema migration queue planning the next ready work item.

## Artifacts

- `docs/token-verifier-record-schema-boundary.md`
- `docs/token-verifier-record-schema-boundary.zh-CN.md`
- `decisions/ADR-0033-token-verifier-record-schema-boundary.md`
- `changes/2026-05-14-define-token-verifier-record-schema-boundary/`

## Open Questions

- Exact migration order for `authentication_device_credentials` and `authentication_access_tokens`.
- Exact repository interface shape for future credential lookup, token issuance, token validation, revocation, and cleanup.
- Exact verifier algorithm, digest format, pepper/secret configuration, and constant-time comparison behavior.
- Whether player and credential relationships become database foreign keys or application-enforced references in the first authentication migration.
- Exact live PostgreSQL verification scope for future authentication migrations and adapters.

## Follow-Up

- Advance W-0076.
- Keep migrations, repositories, adapters, runtime token issuance, validation, logout, cleanup, generated authentication output, Protobuf changes, WebSocket changes, authentication dependencies, and authentication implementation behind later gates.

## Redaction Notes

No secrets, tokens, credential values, account details, provider payloads, authorization headers, cookies, or private data are stored in this conversation log.
