# Conversation: Credential Token Session Schema Gates

Date: 2026-05-14
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-14-define-credential-token-session-schema-gates/`

Related artifacts:

- `docs/credential-token-session-schema-gates.md`
- `docs/credential-token-session-schema-gates.zh-CN.md`
- `decisions/ADR-0029-credential-token-session-schema-gates.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/reference.yaml`
- `.arch/contracts.yaml`
- `.arch/work-items.yaml`

## Context

The maintainer asked the agent to continue advancing the work queue unless a truly necessary decision required human input. W-0071 followed W-0070, which had defined the semantic authentication contract, error, permission, and audit surfaces.

## Maintainer Narrative

The maintainer wants vibit to stay self-bootstrapping and controllable. Routine technical details inside the ratified Go, WebSocket, Protobuf, PostgreSQL, device credential, and opaque token direction may be handled by the agent using professional judgment. True product or architecture branch points should still be brought back for confirmation.

## Agent Response Summary

The agent defined schema gates for credential records, token verifier records, external identity links, runtime session records, audit persistence, and player account lifecycle preservation.

The first posture requires future credential and token verifier schema ratification before implementation. External identity links, runtime sessions, refresh-token storage, session-token vocabulary, WebSocket connection binding, and Redis-like storage remain deferred.

The agent preserved the player account lifecycle schema and did not add migrations, tables, repositories, adapters, runtime lookup code, handlers, routes, Protobuf changes, WebSocket changes, or authentication behavior.

## Decisions

- Define `docs/credential-token-session-schema-gates.md` as the W-0071 standard.
- Record ADR-0029 for credential, token, external identity, session, and audit schema gates.
- Require future credential record schema ratification before `device_credential_login` implementation.
- Require future token verifier schema ratification before opaque access-token validation.
- Keep PostgreSQL as the default durable target for the required future gates.
- Defer external identity storage, runtime session storage, refresh-token storage, session-token vocabulary, WebSocket binding, and Redis-like stores.
- Preserve `player_accounts` and `player_account_events` as lifecycle-only storage.

## Artifacts

- `docs/credential-token-session-schema-gates.md`
- `docs/credential-token-session-schema-gates.zh-CN.md`
- `decisions/ADR-0029-credential-token-session-schema-gates.md`
- `changes/2026-05-14-define-credential-token-session-schema-gates/`

## Open Questions

- Exact credential table shape.
- Exact token verifier table shape.
- Verifier algorithm and comparison rules.
- Whether a future session persistence milestone uses PostgreSQL, a Redis-like store, or a hybrid.
- Whether durable audit persistence becomes a table, event log, outbox, operational log, or another approved store.

## Follow-Up

- Advance W-0072 to add narrow repository checks for selected login/token boundaries.
- Do not implement runtime authentication until a future implementation milestone explicitly authorizes schema ratification, migrations, repositories, adapters, tests, and runtime wiring.

## Redaction Notes

No secrets, tokens, account details, or private data are stored in this conversation log.
