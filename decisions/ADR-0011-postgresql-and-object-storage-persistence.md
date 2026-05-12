# ADR-0011: PostgreSQL And Object Storage Persistence

Status: Accepted
Date: 2026-05-12
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-12-ratify-persistence-direction/`

Related conversations:

- `conversations/2026-05-12-persistence-direction.md`

Related artifacts:

- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `AGENTS.md`
- `README.md`
- `modules/inventory/module.yaml`
- `modules/inventory/AGENTS.md`
- `decisions/ADR-0010-foundational-dependency-adoption.md`

## Context

vibit needs a persistence direction before the first serious Go runtime slice starts. The maintainer proposed setting persistence directly to "psql" and adding MinIO.

The architecture decision should name PostgreSQL rather than `psql`: `psql` is the command-line client, while PostgreSQL is the database engine and persistence platform.

The persistence decision has two different concerns:

- Authoritative durable transactional state for modules and invariants.
- Large object storage for artifacts such as replays, snapshots, exports, binary assets, and diagnostic archives.

PostgreSQL is a mature fit for the first concern. Object storage is useful for the second concern, but it should not be treated as a replacement for relational state or as required for every runtime slice.

MinIO is S3-compatible object storage and can be useful for local and self-hosted deployments. However, the public `minio/minio` GitHub repository is archived as of 2026-04-25 and reports AGPL-3.0 licensing metadata. That makes it a candidate that requires explicit dependency adoption review before becoming a mandatory runtime dependency.

External facts checked:

- PostgreSQL official license page: https://www.postgresql.org/about/licence/
- MinIO GitHub repository: https://github.com/minio/minio

## Decision

PostgreSQL is the first authoritative durable relational store for vibit runtime state.

S3-compatible object storage is the planned object-storage abstraction for large object artifacts.

MinIO is the preferred local/self-hosted candidate for S3-compatible object storage, but it is not a mandatory runtime dependency yet. Making MinIO required requires a dependency adoption record that evaluates licensing, maintenance status, operational fit, replacement path, and the concrete object-storage use case.

Domain modules must not depend directly on PostgreSQL drivers, migration tools, S3 SDKs, or MinIO clients. Platform storage adapters own those dependencies and expose vibit-owned repository and object-storage interfaces to modules.

This ADR does not choose:

- A Go PostgreSQL driver.
- A migration tool.
- A transaction boundary implementation.
- An outbox or event-log storage design.
- An S3 SDK.
- A MinIO deployment model.

Those choices require focused change specs or adoption records under `ADR-0010`.

## Alternatives Considered

- Start with only in-memory persistence for the first slice.
- Use PostgreSQL as the first authoritative durable store.
- Use a document database or key-value store as the first persistence engine.
- Make MinIO a required runtime dependency immediately.
- Define only an S3-compatible object-storage interface and keep MinIO as a candidate.

## Rationale

Game and backend servers need authoritative state, transactional consistency, constraints, migrations, and clear ownership of data. PostgreSQL gives vibit a serious default for that without forcing domain modules to know driver details.

The agent-native benefit is also strong: PostgreSQL schema, migration files, repository interfaces, and transaction boundaries can become explicit artifacts that agents can inspect, verify, and generate against.

Object storage solves a different problem. It is appropriate for large immutable or semi-structured artifacts, but it is a poor primary mechanism for enforcing gameplay invariants or transactional module state.

Using the S3-compatible interface as the abstraction keeps deployments flexible. MinIO remains valuable for local development and self-hosting, but its license and maintenance route should be checked before the framework requires it.

## Agent Reasoning Summary

The correct split is PostgreSQL for authoritative relational state and S3-compatible object storage for large artifacts. This gives future agents a clear persistence model while avoiding premature coupling to a specific object-storage server or SDK.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: high
  implementation_cost: medium
  reversibility: medium
  long_term_maintainability: high
confidence: high
```

## Consequences

- The first persistent module repositories should target PostgreSQL through vibit-owned repository interfaces.
- Migration expectations and transaction boundaries must be declared before persistent repository implementation.
- In-memory repositories may still be used only as tests, fakes, or pre-persistence scaffolding, not as the long-term authoritative store.
- Object storage should not be introduced into a module unless a contract or platform capability has a real large-object use case.
- MinIO can be used later after adoption review, but agents must not add it casually as a required runtime dependency.
- Architecture checks should eventually forbid domain modules from importing PostgreSQL drivers, S3 SDKs, or MinIO clients directly.

## Reversal Conditions

Revisit this decision if early runtime work shows PostgreSQL creates unacceptable operational complexity, weakens agent-readable data ownership, or cannot support the framework's transactional and verification needs.

Revisit the MinIO candidate status if its licensing, maintenance model, or operational fit changes, or if another S3-compatible implementation becomes a better default for local and self-hosted deployments.

## Follow-Up

- Define the first Go repository interface layout before persistent implementation begins.
- Choose the first PostgreSQL driver through a dependency adoption record.
- Define migration file conventions and verification commands.
- Define transaction boundary rules for command handlers and event publication.
- Decide whether an outbox pattern is required for the first persistent event flow.
- Create a separate adoption record before adding any S3 SDK or MinIO dependency.
