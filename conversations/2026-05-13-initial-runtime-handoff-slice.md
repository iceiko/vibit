# Conversation: Initial Runtime Handoff Slice

Date: 2026-05-13
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-13-implement-initial-runtime-handoff-slice/`

Related artifacts:

- `runtime/internal/app/`
- `runtime/internal/platform/protocol/protobuf/`
- `runtime/internal/generated/proto/`
- `.arch/runtime.yaml`
- `tools/vibit`

## Context

The maintainer asked to continue. The local environment now has Go, Buf, and `protoc`, so the project no longer needed to avoid Go runtime code or generated Protobuf output for toolchain reasons.

Previous ADRs had already accepted Go, WebSocket, Protobuf, generated output rules, and the runtime protocol adapter boundary. ADR-0018 specifically said the first implementation should introduce narrow Go handoff types before WebSocket wiring grows.

## Maintainer Narrative

The maintainer wants steady progress toward a serious long-lived system rather than a disposable demo. Earlier guidance emphasized that implementation should happen only when the project has enough preparation to stay self-bootstrapping and controllable.

## Agent Response Summary

The agent chose a narrow first Go runtime slice:

- Generate Go Protobuf output from the existing protocol and inventory `.proto` sources.
- Add pure application handoff types under `runtime/internal/app/`.
- Add Protobuf protocol adapter helpers under `runtime/internal/platform/protocol/protobuf/`.
- Convert generated Protobuf envelopes into application route requests.
- Decode generated payloads in the Protobuf adapter instead of leaking generated Protobuf types into application dispatch.
- Add runtime import-boundary checks so foundational dependencies remain inside their declared owner packages.

The agent intentionally did not implement WebSocket transport, PostgreSQL persistence, migrations, authentication, or full application dispatch in this slice.

## Decisions

- Start Go runtime implementation with generated Protobuf output and protocol handoff helpers.
- Keep generated Protobuf imports out of `runtime/internal/app/`.
- Let the Protobuf adapter own generated envelope and payload decoding.
- Add import-boundary checks now that Go source files exist.

## Artifacts

- Generated `runtime/internal/generated/proto/vibit/protocol/v1/envelope.pb.go`.
- Generated `runtime/internal/generated/proto/vibit/inventory/v1/inventory.pb.go`.
- Added application handoff types under `runtime/internal/app/`.
- Added Protobuf envelope adapter helpers and tests under `runtime/internal/platform/protocol/protobuf/`.
- Updated runtime manifests, guides, README files, runtime adapter docs, rule catalog, and checks.

## Open Questions

- Should the first persistent inventory repository start directly with PostgreSQL, or should a test-only in-memory fake be introduced behind the same interface?
- Should route registration be generated from contracts before the first application dispatcher implementation?
- Should payload type naming remain Protobuf full message name, or should vibit add an explicit registry artifact before public clients exist?

## Follow-Up

- Add application dispatch tests before domain runtime behavior grows.
- Decide the first repository implementation shape before persistence code begins.
- Add WebSocket transport only after protocol adapter behavior is covered by tests.

## Redaction Notes

No secrets, tokens, account identifiers, or private data were included in this conversation log.
