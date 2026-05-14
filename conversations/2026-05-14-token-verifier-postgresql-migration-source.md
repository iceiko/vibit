# Conversation: Token Verifier PostgreSQL Migration Source

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-add-token-verifier-postgresql-migration-source/`

Related artifacts:

- `runtime/migrations/postgres/000004_create_authentication_access_tokens.sql`
- `docs/token-verifier-record-schema-boundary.md`
- `docs/authentication-schema-migration-queue.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`

## Context

The maintainer asked the agent to continue advancing bounded work items unless a true decision or blocker required confirmation. W-0078 follows W-0077 credential migration source creation.

## Maintainer Narrative

The maintainer wants vibit to advance through professional judgment inside accepted architecture decisions, with Nakama and Pitaya remaining active references for game-backend capability coverage while vibit's agent-native boundaries remain authoritative.

## Agent Response Summary

The agent added the SQL-first PostgreSQL migration source for the ratified `authentication_access_tokens` schema.

The migration creates only the token verifier table, uses the `runtime.authentication` module trace, stores non-plaintext token lookup and verifier digests, links each token verifier to a player account and credential record, and records expiration, revocation, rotation lineage, retention, and cleanup eligibility semantics without modifying player lifecycle or credential tables.

The agent updated manifests, standards, guides, and checks to record that both credential and token verifier migration sources exist while repository interfaces, PostgreSQL adapters, runtime credential lookup, token issuance, token validation, logout, refresh, cleanup, generated authentication shapes, Protobuf messages, WebSocket proof carriers, WebSocket handshake authentication, authentication dependencies, and production authentication behavior remain deferred.

## Decisions

- Add `runtime/migrations/postgres/000004_create_authentication_access_tokens.sql`.
- Use `-- Module: runtime.authentication` as the owning module trace.
- Use `authentication_access_tokens` as the only new table in this step.
- Store `token_lookup_digest` and `token_verifier_digest` as `BYTEA`.
- Reference `player_accounts(player_id)` and `authentication_device_credentials(credential_record_id)`.
- Keep authentication repository interfaces and PostgreSQL adapters for later work.

## Artifacts

- `runtime/migrations/postgres/000004_create_authentication_access_tokens.sql`
- `changes/2026-05-14-add-token-verifier-postgresql-migration-source/`

## Open Questions

- Exact authentication migration static checks for W-0079.
- Exact storage-neutral authentication repository shape for W-0080.
- Exact PostgreSQL adapter boundary for W-0081.
- Exact verifier algorithm and secret configuration boundary for future runtime implementation.

## Follow-Up

- Advance W-0079.
- Keep runtime authentication, token validation, logout, cleanup, generated authentication output, Protobuf changes, WebSocket changes, authentication dependencies, and production authentication behavior behind later gates.

## Redaction Notes

No secrets, raw token values, raw credential values, credential digests, token digests, account details, provider payloads, authorization headers, cookies, or private data are stored in this conversation log.
