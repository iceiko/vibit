# Conversation: Authentication Schema Migration Queue

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-plan-authentication-schema-migration-queue/`

Related artifacts:

- `docs/authentication-schema-migration-queue.md`
- `docs/authentication-schema-migration-queue.zh-CN.md`
- `decisions/ADR-0034-authentication-schema-migration-queue.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`

## Context

The maintainer asked the agent to continue advancing bounded work items unless a true decision or blocker required confirmation. W-0076 follows credential and token verifier schema-boundary ratification in M-014.

## Maintainer Narrative

The maintainer wants vibit to stay self-bootstrapping and controlled. Preparation is valuable when it prevents future drift, but routine technical sequencing inside an accepted direction should be handled by the agent's professional judgment.

## Agent Response Summary

The agent planned the authentication schema migration queue without adding schema or runtime behavior.

The queue places credential migration before token verifier migration, adds a static migration-check gate before repository interfaces and adapters, and keeps runtime authentication implementation outside M-014.

The agent did not add migration source, tables, repository interfaces, PostgreSQL adapters, runtime credential lookup, token issuance, token validation, logout, refresh, cleanup jobs, generated authentication shapes, Protobuf messages, WebSocket proof carriers, WebSocket handshake authentication, authentication dependencies, or production authentication behavior.

## Decisions

- Use `W-0077` for the credential PostgreSQL migration source.
- Use `W-0078` for the token verifier PostgreSQL migration source.
- Use `W-0079` for authentication migration static checks.
- Use `W-0080` for the authentication repository interface boundary.
- Use `W-0081` for the authentication PostgreSQL adapter boundary.
- Use `W-0082` for M-014 closeout.
- Plan `000003_create_authentication_device_credentials.sql` before `000004_create_authentication_access_tokens.sql`.
- Keep runtime authentication implementation outside M-014.

## Artifacts

- `docs/authentication-schema-migration-queue.md`
- `docs/authentication-schema-migration-queue.zh-CN.md`
- `decisions/ADR-0034-authentication-schema-migration-queue.md`
- `changes/2026-05-14-plan-authentication-schema-migration-queue/`

## Open Questions

- Exact SQL column types and constraints for `authentication_device_credentials` in W-0077.
- Exact SQL column types and constraints for `authentication_access_tokens` in W-0078.
- Exact static migration check implementation in W-0079.
- Exact storage-neutral authentication repository shape in W-0080.
- Whether M-014 closeout opens a confirmation gate or a concrete authentication implementation milestone.

## Follow-Up

- Advance W-0077.
- Keep runtime authentication, token validation, logout, cleanup, generated authentication output, Protobuf changes, WebSocket changes, authentication dependencies, and production authentication behavior behind later gates.

## Redaction Notes

No secrets, tokens, credential values, account details, provider payloads, authorization headers, cookies, or private data are stored in this conversation log.
