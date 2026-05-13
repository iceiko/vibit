# Conversation: Runtime Protocol Adapter Boundary

Date: 2026-05-13

Related changes:

- `changes/2026-05-13-define-runtime-protocol-adapter-boundary/`

Related artifacts:

- `docs/runtime-protocol-adapter.md`
- `docs/runtime-protocol-adapter.zh-CN.md`
- `decisions/ADR-0018-runtime-protocol-adapter-boundary.md`
- `.arch/runtime.yaml`
- `.arch/protocol.yaml`
- `tools/vibit`

## Context

The maintainer asked to continue development. The previous generated output standard was committed locally, but GitHub push attempts failed because the environment could not connect to GitHub on port 443.

The local environment was checked for Go, Buf, and `protoc`; none were available.

## Maintainer Narrative

The maintainer wants steady forward progress without rushing into a shallow demo. Previous project direction emphasizes making the project self-bootstrapping, controllable, and explicitly designed for future agents.

## Agent Response Summary

The agent chose not to fake generated output or Go runtime tests without the required toolchain.

Instead, the agent defined the runtime protocol adapter boundary so future WebSocket, Protobuf, application dispatch, and domain runtime code has a clear handoff model before implementation starts.

## Decisions

- Add a runtime protocol adapter boundary standard before Go runtime implementation.
- Keep WebSocket transport, Protobuf protocol adaptation, application dispatch, domain logic, and generated output responsibilities separate.
- Extend runtime checks to verify that the boundary standard is wired into manifests and agent guidance.

## Artifacts

- `docs/runtime-protocol-adapter.md`
- `docs/runtime-protocol-adapter.zh-CN.md`
- `decisions/ADR-0018-runtime-protocol-adapter-boundary.md`
- `changes/2026-05-13-define-runtime-protocol-adapter-boundary/`

## Open Questions

- When will the local environment have Go, Buf, and Protobuf tooling available?
- Should the first Go implementation start with handoff types and tests before generated Protobuf output exists?

## Follow-Up

- Verify the new boundary standard and runtime checks.
- Commit the change locally.
- Attempt to push when GitHub network connectivity is available.

## Redaction Notes

No secrets were included.
