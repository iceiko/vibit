# Conversation: Credential PostgreSQL Migration Source

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-add-credential-postgresql-migration-source/`

Related artifacts:

- `runtime/migrations/postgres/000003_create_authentication_device_credentials.sql`
- `docs/credential-record-schema-boundary.md`
- `docs/authentication-schema-migration-queue.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`

## Context

The maintainer asked the agent to continue advancing bounded work items unless a true decision or blocker required confirmation. W-0077 follows W-0076 authentication schema migration queue planning.

## Maintainer Narrative

The maintainer wants vibit to stay self-bootstrapping, controlled, and agent-native. Routine technical sequencing inside an accepted architecture direction should be handled by the agent's professional judgment.

## Agent Response Summary

The agent added the SQL-first PostgreSQL migration source for the ratified `authentication_device_credentials` schema.

The migration creates only the credential verifier table, uses the `runtime.authentication` module trace, stores non-plaintext lookup and verifier digests, preserves one active device credential per player for the first posture, and references `player_accounts(player_id)` without modifying player account lifecycle tables.

The agent updated manifests, standards, guides, and checks to record that the credential migration source exists while token verifier migration, repository interfaces, PostgreSQL adapters, runtime credential lookup, token issuance, token validation, logout, refresh, cleanup, generated authentication shapes, Protobuf messages, WebSocket proof carriers, WebSocket handshake authentication, authentication dependencies, and production authentication behavior remain deferred.

## Decisions

- Add `runtime/migrations/postgres/000003_create_authentication_device_credentials.sql`.
- Use `-- Module: runtime.authentication` as the owning module trace.
- Use `authentication_device_credentials` as the only new table in this step.
- Store `credential_lookup_digest` and `credential_verifier_digest` as `BYTEA`.
- Use a partial unique index to enforce at most one active device credential per player in the first posture.
- Keep token verifier storage for W-0078.

## Artifacts

- `runtime/migrations/postgres/000003_create_authentication_device_credentials.sql`
- `changes/2026-05-14-add-credential-postgresql-migration-source/`

## Open Questions

- Exact token verifier SQL shape for W-0078.
- Exact authentication migration static checks for W-0079.
- Exact storage-neutral authentication repository shape for W-0080.
- Exact verifier algorithm and secret configuration boundary for future runtime implementation.

## Follow-Up

- Advance W-0078.
- Keep runtime authentication, token validation, logout, cleanup, generated authentication output, Protobuf changes, WebSocket changes, authentication dependencies, and production authentication behavior behind later gates.

## Redaction Notes

No secrets, raw credential values, raw tokens, credential digests, token digests, account details, provider payloads, authorization headers, cookies, or private data are stored in this conversation log.
