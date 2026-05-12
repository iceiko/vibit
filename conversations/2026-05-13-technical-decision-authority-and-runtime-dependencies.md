# Conversation: Technical Decision Authority And Runtime Dependencies

Date: 2026-05-13
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-13-adopt-first-runtime-dependencies/`

Related artifacts:

- `.arch/dependencies.yaml`
- `.arch/runtime.yaml`
- `decisions/ADR-0012-agent-technical-decision-authority.md`
- `decisions/ADR-0013-first-go-runtime-dependencies.md`

## Context

The previous dependency adoption standard required the agent to discuss branch-producing foundational dependency choices before accepting concrete packages. The maintainer then clarified that some technical branch points should be decided by professional evaluation rather than repeatedly blocking on confirmation.

## Maintainer Narrative

The maintainer said:

```text
这些你按照专业评估来定。这些问题越专业越好，你自己评估就可以。以后这种问题就是要来跟我确认，我允许你自己评估决定，你才能自己评估决定。
```

This means the agent may evaluate and decide technical dependency and runtime preparation questions when the maintainer grants that authority.

It does not remove the need to ask before product direction changes, constitutional changes, major architecture pattern changes, scope expansion, cost or operations commitments, licensing-risk acceptance, or decisions that redefine module ownership.

## Agent Response Summary

The agent interpreted the authorization as bounded technical delegation. It proceeded to evaluate first foundational Go runtime dependencies using maintenance, license, ecosystem, agent-readability, boundary, replacement, and verification criteria.

The agent accepted mature dependencies only where the boundary is narrow and the first runtime slice will need the decision:

- `github.com/coder/websocket` for the platform WebSocket transport adapter.
- `google.golang.org/protobuf` and `protoc-gen-go` for Go Protobuf runtime and generated Go bindings.
- `buf.build` CLI for Protobuf linting, breaking checks, formatting, and generation orchestration.
- `github.com/jackc/pgx/v5` for PostgreSQL platform persistence adapters.
- `github.com/pressly/goose/v3` for SQL-first migration tooling.

The agent left S3 client, MinIO deployment, observability, and external test framework adoption deferred or proposed because the first runtime slice does not yet require them.

## Decisions

- Record bounded agent technical decision authority in ADR-0012.
- Record first accepted Go runtime dependencies in ADR-0013.
- Keep accepted dependencies behind platform adapters or generation tooling.
- Continue asking the maintainer before decisions outside delegated technical evaluation.

## Artifacts

- Added `decisions/ADR-0012-agent-technical-decision-authority.md`.
- Added `decisions/ADR-0013-first-go-runtime-dependencies.md`.
- Updated `.arch/dependencies.yaml`.
- Updated `.arch/runtime.yaml`.
- Updated repository agent guides and architecture README files.
- Added `changes/2026-05-13-adopt-first-runtime-dependencies/`.

## Open Questions

- What exact Go package layout should hold platform adapters and generated protocol packages?
- What transaction boundary convention should command handlers use once persistence implementation begins?
- What manifest-to-Protobuf alignment checker should be implemented before broad protocol generation?
- When is the first object-storage use case concrete enough to adopt an S3 SDK or MinIO deployment?

## Follow-Up

- Define the Go package layout and adapter boundaries before creating `go.mod`.
- Define repository, migration, and transaction boundary conventions before implementing persistent repositories.
- Generate the first Protobuf wire schema only after the semantic contract to Protobuf alignment path is declared.

## Redaction Notes

No secrets, tokens, account identifiers, or private data were included in this conversation log.
