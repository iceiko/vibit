# Request

## Original Request

```text
继续
```

## Clarified Requirement

Now that Go, Buf, and `protoc` are available locally, start the first narrow Go runtime implementation slice.

This slice should:

- Generate Go Protobuf output from the existing `.proto` sources.
- Add a small, agent-readable runtime protocol adapter package.
- Introduce the first narrow Go handoff types and tests.
- Keep WebSocket transport, PostgreSQL persistence, outbox behavior, and broader application wiring out of scope for now.

## User-Visible Outcome

The repository should now contain the first real Go runtime code path, not only manifests and skeleton directories.

Future agents should be able to inspect the runtime package, the generated Protobuf output, and the tests to understand how the first protocol handoff is intended to work.

## Non-Goals

- Do not implement a WebSocket server.
- Do not implement PostgreSQL repositories or migrations.
- Do not add an outbox or event delivery system.
- Do not add authentication or session lifecycle behavior.
- Do not change the public command, query, event, error, or permission contracts.
- Do not introduce a broad application dispatcher yet.

## Unknowns

- The exact first repository implementation shape remains open.
- The first transport wiring remains deferred.
- The boundary between protocol adapter helpers and application dispatch may need refinement after the first slice exists.

## Acceptance Criteria

- Generated Go Protobuf output exists under `runtime/internal/generated/proto/`.
- The runtime module includes a small protocol adapter package with narrow handoff helpers.
- The first Go runtime tests compile and pass locally.
- `node tools/vibit check runtime` now exercises Go source and test files.
- `node tools/vibit check generated`, `node tools/vibit check runtime`, and `node tools/vibit check all` pass.
