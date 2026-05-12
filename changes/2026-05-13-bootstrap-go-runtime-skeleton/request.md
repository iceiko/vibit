# Request

## Original Request

```text
继续推进
```

## Clarified Requirement

Create the first Go runtime skeleton according to ADR-0014 before adding business runtime code.

The skeleton should establish the Go module root, package boundary directories, runtime agent guide, Protobuf source root, and runtime verification semantics.

## User-Visible Outcome

Agents can now see the concrete runtime and Protobuf roots on disk instead of inferring them from manifests alone.

## Non-Goals

- Do not implement Inventory behavior.
- Do not add WebSocket, Protobuf, PostgreSQL, goose, or S3 dependencies to `go.mod`.
- Do not create generated Go output.
- Do not create `.proto` message files.
- Do not create database migrations.

## Unknowns

- The local environment currently does not have the `go` command installed.
- The first repository implementation shape remains open: fake first, PostgreSQL first, or both.
- The manifest-to-Protobuf alignment checker is still not implemented.

## Acceptance Criteria

- `runtime/go.mod` exists with the ratified module path.
- Runtime package boundary directories from ADR-0014 are present.
- Runtime and Protobuf source roots have agent-readable guides.
- `node tools/vibit check runtime` reads `runtime/go.mod`, not repository-root `go.mod`.
- Runtime checks pass for a skeleton with no Go source files.
- Existing architecture, schema, memory, contract, generated, runtime, and all checks pass.
