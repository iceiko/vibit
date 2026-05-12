# Request

Date: 2026-05-12
Change ID: `ratify-persistence-direction`
Type: standard

## Maintainer Request

The maintainer asked to continue architecture preparation and proposed a persistence direction:

```text
继续，另外，我们的持久化目前直接定成psql，然后再加一个minIO，你觉得怎么样？
```

## Clarified Requirement

Evaluate whether PostgreSQL and MinIO should be ratified as the persistence direction, then record the bounded decision in architecture artifacts if appropriate.

Use PostgreSQL rather than `psql` as the architecture term. `psql` is a client command; PostgreSQL is the database engine and persistence platform.

## User-Visible Outcome

The repository should clearly state:

- PostgreSQL is the first authoritative durable relational store.
- S3-compatible object storage is the planned abstraction for large artifacts.
- MinIO is the preferred local/self-hosted candidate, but not a mandatory runtime dependency until dependency adoption review is complete.
- Domain modules must depend on vibit-owned storage interfaces rather than direct database or object-storage clients.

## Non-Goals

- Do not implement runtime persistence code.
- Do not choose a PostgreSQL Go driver yet.
- Do not choose a migration tool yet.
- Do not choose an S3 SDK yet.
- Do not add MinIO as a dependency or deployment requirement yet.
- Do not modify historical change specs that described earlier in-memory assumptions.

## Acceptance Criteria

- [x] ADR records the persistence decision and alternatives.
- [x] `.arch/runtime.yaml` records PostgreSQL and S3-compatible object-storage direction.
- [x] `.arch/conventions.yaml` links the persistence decision.
- [x] Repository and module agent guides warn against direct persistence client dependencies inside domain modules.
- [x] README files mention the bounded persistence direction.
- [x] Conversation memory records the maintainer's original proposal with no secrets.
