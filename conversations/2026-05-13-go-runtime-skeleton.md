# Conversation: Go Runtime Skeleton

Date: 2026-05-13
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-13-bootstrap-go-runtime-skeleton/`

Related artifacts:

- `runtime/go.mod`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `proto/README.md`
- `proto/README.zh-CN.md`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

After ADR-0014 defined the Go runtime layout, the next step was to make the layout concrete without starting business implementation.

The existing `check runtime` command still assumed a repository-root `go.mod`, which no longer matched ADR-0014.

## Maintainer Narrative

The maintainer asked:

```text
继续推进
```

The maintainer had previously authorized the agent to professionally evaluate and decide technical sub-decisions inside already ratified directions.

## Agent Response Summary

The agent proceeded with a bounded runtime skeleton:

- Create `runtime/go.mod` with module path `github.com/iceiko/vibit/runtime`.
- Create ADR-0014 package boundary directories with `.gitkeep` placeholders.
- Create `runtime/AGENTS.md` and `runtime/AGENTS.zh-CN.md`.
- Create `proto/README.md`, `proto/README.zh-CN.md`, and `proto/vibit/inventory/v1/`.
- Update `tools/vibit check runtime` to check `runtime/go.mod` and skeleton paths.
- Keep Go business code, dependencies, generated files, Protobuf files, migrations, and runtime tests deferred.

The local environment did not provide the `go` command during this change, so Go test verification was not claimed.

## Decisions

- Treat `runtime/go.mod` plus ADR-0014 directories as a skeleton, not business runtime implementation.
- Runtime checks should pass for a skeleton with no Go source files.
- Once Go source files exist, runtime checks should require Go test files and a local Go toolchain.

## Artifacts

- Added the Go runtime skeleton under `runtime/`.
- Added the Protobuf source-root skeleton under `proto/`.
- Updated runtime manifests and repository guides.
- Updated `tools/vibit` runtime check semantics.
- Updated `rules/check-rules.json` with runtime skeleton rule metadata.

## Open Questions

- Should the first repository implementation start with an in-memory fake, PostgreSQL, or both?
- Should the manifest-to-Protobuf alignment checker come before the first `.proto` file?
- Should import-boundary checks be implemented before the first Go source file that imports third-party packages?

## Follow-Up

- Define the manifest-to-Protobuf alignment check before broad protocol generation.
- Add import-boundary checks once Go source files exist.
- Begin the first Inventory runtime slice only after the source, generated, handwritten, test, and verification path is clear.

## Redaction Notes

No secrets, tokens, account identifiers, or private data were included in this conversation log.
