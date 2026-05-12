# Conversation: Go Runtime Layout And Boundaries

Date: 2026-05-13
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-13-define-go-runtime-layout/`

Related artifacts:

- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `decisions/ADR-0014-go-runtime-layout-and-boundaries.md`

## Context

After accepting the first foundational Go runtime dependencies, the remaining runtime readiness questions were package layout, platform adapter boundaries, Protobuf/Buf layout, migration layout, and transaction boundary rules.

These questions must be answered before creating `go.mod` or adding Go runtime code, otherwise agents would need to infer structural choices from partial implementation.

## Maintainer Narrative

The maintainer asked:

```text
继续推进
```

The maintainer had previously authorized the agent to professionally evaluate and decide technical sub-decisions inside already ratified directions.

## Agent Response Summary

The agent proceeded with a conservative runtime layout decision:

- Use `runtime/` as the first Go module root.
- Keep `.proto` source files under repository-root `proto/`.
- Generate Go protocol code under `runtime/internal/generated/proto/`.
- Keep domain module logic under `runtime/internal/modules/<module>/`.
- Keep third-party transport and persistence dependencies inside `runtime/internal/platform/`.
- Keep SQL-first PostgreSQL migrations under `runtime/migrations/postgres/`.
- Route state-changing commands through an application transaction boundary before repository mutation and event publication.

No Go code, `go.mod`, Protobuf output, or migrations were added in this change.

## Decisions

- Add ADR-0014 for Go runtime layout and boundaries.
- Update `.arch/runtime.yaml` with machine-readable layout and boundary rules.
- Keep implementation code deferred until the package layout is recorded.

## Artifacts

- Added `decisions/ADR-0014-go-runtime-layout-and-boundaries.md`.
- Updated `.arch/runtime.yaml`.
- Updated `.arch/conventions.yaml`.
- Updated AGENTS and README pairs.
- Added `changes/2026-05-13-define-go-runtime-layout/`.

## Open Questions

- Should the first runtime implementation include an in-memory repository fake before PostgreSQL integration?
- Should the first event publication path include an outbox table immediately, or defer outbox until cross-boundary event delivery exists?
- Should `tools/vibit` generate the initial Go package skeleton before or after `go.mod` is created?

## Follow-Up

- Create the first Go module and skeletal runtime packages according to ADR-0014.
- Add import-boundary checks once Go files exist.
- Define the manifest-to-Protobuf alignment checker before broad protocol generation.

## Redaction Notes

No secrets, tokens, account identifiers, or private data were included in this conversation log.
