# Conversation: Persistence Direction

Date: 2026-05-12
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-12-ratify-persistence-direction/`

Related artifacts:

- `decisions/ADR-0011-postgresql-and-object-storage-persistence.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `modules/inventory/module.yaml`
- `modules/inventory/AGENTS.md`
- `modules/inventory/AGENTS.zh-CN.md`

## Context

After ratifying Go, WebSocket, and Protobuf as the first runtime direction, the maintainer asked to continue and raised the persistence stack before runtime implementation began.

## Maintainer Narrative

The maintainer asked:

```text
继续，另外，我们的持久化目前直接定成psql，然后再加一个minIO，你觉得怎么样？
```

This continues the maintainer's earlier requirement that major server architecture choices be discussed and analyzed before being ratified.

## Agent Response Summary

The agent recommended splitting the decision into two different persistence concerns:

- PostgreSQL, not `psql`, should be the architecture term for the first authoritative durable relational store.
- S3-compatible object storage should be the abstraction for large artifacts.
- MinIO can be the preferred local/self-hosted candidate, but should not become a mandatory runtime dependency without a dependency adoption record.

The agent checked current public metadata and noted that the public `minio/minio` GitHub repository is archived as of 2026-04-25 and reports AGPL-3.0 license metadata. This supports treating MinIO as a candidate requiring explicit adoption review rather than a casually required dependency.

External facts checked:

- PostgreSQL official license page: https://www.postgresql.org/about/licence/
- MinIO GitHub repository: https://github.com/minio/minio

## Decisions

- PostgreSQL is the first authoritative durable relational store for vibit runtime state.
- S3-compatible object storage is the planned abstraction for large object artifacts.
- MinIO is the preferred local/self-hosted candidate for S3-compatible object storage, pending dependency adoption before required runtime use.
- Domain modules must use vibit-owned repository and storage interfaces rather than depending directly on PostgreSQL drivers, S3 SDKs, or MinIO clients.

## Artifacts

- Added `decisions/ADR-0011-postgresql-and-object-storage-persistence.md`.
- Added `changes/2026-05-12-ratify-persistence-direction/`.
- Updated `.arch/runtime.yaml`.
- Updated `.arch/conventions.yaml`.
- Updated README and AGENTS files in English and Simplified Chinese.
- Updated inventory module metadata and agent guides.

## Open Questions

- Which Go PostgreSQL driver should vibit adopt first?
- Which migration tool or migration file convention should vibit use?
- Where should transaction boundaries live in generated and handwritten runtime shape?
- Does the first persistent event flow require an outbox pattern?
- Which S3-compatible SDK should be adopted when object storage is first needed?

## Follow-Up

- Create a dependency adoption record before selecting a PostgreSQL driver.
- Define repository interface layout and transaction boundary conventions before implementing persistent repositories.
- Create a separate adoption record before adding MinIO or any S3 SDK as a dependency.

## Redaction Notes

No secrets, tokens, account identifiers, or private data were included in this conversation log.
