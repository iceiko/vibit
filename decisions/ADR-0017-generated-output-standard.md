# ADR-0017: Generated Output Standard

Status: Accepted
Date: 2026-05-13
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-13-add-generated-output-standard/`

Related conversations:

- `conversations/2026-05-13-generated-output-standard.md`

Related artifacts:

- `docs/generated-output.md`
- `docs/generated-output.zh-CN.md`
- `.arch/runtime.yaml`
- `.arch/protocol.yaml`
- `runtime/internal/generated/proto/`
- `buf.gen.yaml`
- `tools/vibit`

## Context

vibit has accepted Protobuf as the first wire schema format and Buf plus `protoc-gen-go` as the first generation path for Go Protobuf output.

The repository now has `.proto` sources and a planned generated output root, but generated Go Protobuf files have not been committed. This is the right moment to define the guardrails for generated output before agents start adding generated files or runtime protocol adapters.

The long-term risk is that an agent may fix a protocol or runtime issue by editing generated output directly. That would break the source-before-output model and make future regeneration destructive.

## Decision

Define `docs/generated-output.md` as the generated output standard and add a paired Simplified Chinese translation.

For Go Protobuf output:

- Source files live under `proto/`.
- Generation configuration lives in `buf.yaml` and `buf.gen.yaml`.
- Output lives under `runtime/internal/generated/proto/`.
- Generated Go Protobuf files must use the `*.pb.go` suffix.
- Generated Go Protobuf files must contain the `protoc-gen-go` generated-code marker.
- Generated Go Protobuf files must contain a `// source: <path>.proto` trace that resolves to an existing source under `proto/`.
- Handwritten runtime code is forbidden under `runtime/internal/generated/proto/`.

Extend `node tools/vibit check generated` to verify the current Go Protobuf output root. Empty planned output directories with `.gitkeep` remain valid until generation is actually run.

## Alternatives Considered

- Defer generated output rules until after `buf generate`.
- Allow generated Protobuf output without additional repository checks.
- Require every generated file to include custom vibit trace comments.

## Rationale

Defining the rule before generation prevents a common agent failure mode: fixing generated code instead of the source that produces it.

The check stays pragmatic. It does not require generated output to exist before the toolchain exists, but once files appear under the generated Protobuf root, it verifies they look like `protoc-gen-go` output and trace back to actual `.proto` files.

Using the standard `protoc-gen-go` header avoids inventing custom comments for third-party generated files.

## Agent Reasoning Summary

The repository is moving from standards into generated runtime shape. A small generated-output standard gives future agents a clear rule before they touch generated code, without forcing runtime implementation to start prematurely.

## Decision Weights

```yaml
decision_weights:
  agent_context: high
  human_ergonomics: medium
  implementation_cost: low
  reversibility: high
  long_term_maintainability: high
confidence: high
```

## Consequences

- Future generated Protobuf files must pass `node tools/vibit check generated`.
- Agents have a clear standard for when generated output may be committed.
- Handwritten protocol adapter code remains outside `runtime/internal/generated/proto/`.
- The first generation run can be verified without relying on maintainer memory.

## Reversal Conditions

Revisit this decision if the Go Protobuf generator changes its header format, if Buf output paths change, or if vibit adopts a different generated output layout.

## Follow-Up

- Run `buf lint` and `buf generate` when Buf and the Go Protobuf generator are available.
- Add generated output drift checks once generated files are committed.
- Add import-boundary checks once Go runtime code exists.
