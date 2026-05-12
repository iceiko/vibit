# Conversation: Protocol Alignment Check

Date: 2026-05-13
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-13-add-protocol-alignment-check/`

Related artifacts:

- `tools/vibit`
- `rules/check-rules.json`
- `.arch/contracts.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `proto/README.md`
- `proto/README.zh-CN.md`

## Context

The project has ratified Go as the first server runtime, WebSocket as the first gameplay/client protocol, and Protobuf as the first wire message format. The repository has semantic contracts for the `inventory` proof slice, but no `.proto` source files or generated Protobuf output yet.

The next step was to keep protocol work inside the agent-native workflow before protocol files exist.

## Maintainer Narrative

```text
继续下一步
```

Earlier maintainer direction remains relevant: this is a long-term maintained framework, not a disposable demo; professional technical sub-decisions may be made inside ratified directions, while branch-point decisions should be discussed.

## Agent Response Summary

The agent added a bounded tooling change instead of creating `.proto` files immediately.

The change introduces `node tools/vibit check protocol`, which derives expected Protobuf paths, package names, message names, and field names from registered command, query, and event contracts. Missing `.proto` files are allowed while protocol generation has not started, but existing `.proto` files must align with the semantic contracts.

## Decisions

- Add a manifest-to-Protobuf alignment check before broad Protobuf source creation or generation.
- Keep `contracts/` as the semantic source of truth and `proto/` as the future wire schema source root.
- Do not add Buf configuration, `.proto` files, generated Go Protobuf output, or Go business code in this change.

## Artifacts

- `node tools/vibit check protocol`
- `node tools/vibit check protocol --json`
- Protocol rule metadata under `rules/check-rules.json`
- Architecture alignment status under `.arch/contracts.yaml`, `.arch/runtime.yaml`, and `.arch/conventions.yaml`

## Open Questions

- The exact WebSocket envelope and service shape remains deferred.
- The first repository implementation still needs a decision about fake repository, PostgreSQL repository, or both.
- Event delivery and outbox storage remain deferred.

## Follow-Up

- When `.proto` files are first created, run `node tools/vibit check protocol` before and after the change.
- Later generation tooling should reuse the same manifest-to-Protobuf alignment expectations instead of inventing a separate mapping.

## Redaction Notes

No secrets, tokens, account details, or private data are included.
