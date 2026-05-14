# ADR-0034: Authentication Schema Migration Queue

Status: Accepted
Date: 2026-05-14
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-14-plan-authentication-schema-migration-queue/`

Related conversations:

- `conversations/2026-05-14-authentication-schema-migration-queue.md`

Related artifacts:

- `docs/authentication-schema-migration-queue.md`
- `docs/authentication-schema-migration-queue.zh-CN.md`
- `docs/credential-record-schema-boundary.md`
- `docs/token-verifier-record-schema-boundary.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`

## Context

W-0074 ratified the credential record schema boundary for `authentication_device_credentials`. W-0075 ratified the token verifier record schema boundary for `authentication_access_tokens`.

Both boundaries are schema-only decisions. They do not create SQL migration source, repository interfaces, PostgreSQL adapters, runtime authentication behavior, token validation, WebSocket proof carriers, Protobuf messages, or authentication dependencies.

The next risk is that an agent may jump from schema ratification directly to a broad implementation. W-0076 creates a bounded queue so the next steps remain self-bootstrapping and verifiable.

## Decision

Plan the authentication schema migration queue for the rest of M-014.

The planned work order is:

```text
W-0077 Add credential PostgreSQL migration source
W-0078 Add token verifier PostgreSQL migration source
W-0079 Add authentication migration static checks
W-0080 Define authentication repository interface boundary
W-0081 Define authentication PostgreSQL adapter boundary
W-0082 Close credential and token verifier schema ratification milestone
```

The first authentication migration order is:

```text
000003_create_authentication_device_credentials.sql
000004_create_authentication_access_tokens.sql
```

Credential migration comes first because token verifier records need credential linkage. Token verifier migration comes second because token lifecycle, rotation, revocation, and cleanup semantics depend on that linkage.

Repository interfaces and PostgreSQL adapter boundaries remain later steps. Runtime authentication implementation remains outside M-014.

This decision does not add migration source, tables, repository interfaces, PostgreSQL adapters, runtime credential lookup, token issuance, token validation, logout, refresh, cleanup, generated authentication shapes, Protobuf messages, WebSocket proof carriers, WebSocket handshake authentication, authentication dependencies, or production authentication behavior.

## Alternatives Considered

- Add both authentication migrations in one change.
- Add credential and token migrations, repositories, and adapters together.
- Start runtime authentication after schema boundaries without migration planning.
- Close M-014 immediately after W-0075.
- Add repository interfaces before migration sources.
- Add adapters before static migration checks.
- Ask for a maintainer decision before planning the queue.

## Rationale

The project is agent-native, so the sequence matters as much as the destination. A queue lets later agents continue with a narrow unit of work instead of reconstructing the plan from memory.

Credential migration before token migration preserves the causal dependency between device credentials and access-token rotation. Static checks after the first migrations are useful because they turn the newly allowed authentication schema into enforceable boundaries before repositories and adapters start using it.

Repository interfaces after migrations keep storage-neutral code grounded in actual schema boundaries. Adapter boundary before implementation follows the pattern already used for player account persistence and keeps SQL behavior reviewable before runtime authentication depends on it.

Closing M-014 after adapter boundary planning creates a natural stop before production authentication implementation, which may require a separate confirmation or implementation milestone.

## Agent Reasoning Summary

The correct next step is not to implement authentication. The correct next step is to turn the now-ratified credential and token verifier schema decisions into a deterministic queue that future agents can execute one bounded step at a time.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  schema_first_discipline: high
  migration_reviewability: high
  implementation_deferral: high
  security_boundary_clarity: high
  similarity_to_existing_player_persistence_pattern: high
  game_backend_reference_alignment: medium
  delivery_speed: medium
  long_term_maintainability: high
confidence: high
```

## Consequences

- W-0077 becomes the next ready work item.
- Authentication schema work proceeds in small, ordered steps.
- Credential and token verifier migrations remain separate.
- Repository interfaces remain after migration source and static checks.
- PostgreSQL adapter implementation remains after interface and adapter boundary steps.
- Runtime authentication remains blocked until a later milestone explicitly authorizes it.

## Reversal Conditions

Revisit this decision if:

- Migration sequence numbers are already occupied by other committed migration files.
- A security review requires authentication static checks before any authentication migration source is added.
- A future distributed runtime decision changes the default durable authentication target away from PostgreSQL.
- The maintainer explicitly asks to pause schema migration work and select a different milestone direction.
- A future session persistence decision requires token verifier storage to be redesigned before migration source is created.

## Follow-Up

- Advance W-0077.
- Keep runtime authentication, token validation, logout, cleanup, generated authentication output, Protobuf changes, WebSocket changes, authentication dependencies, and production authentication behavior behind later gates.
