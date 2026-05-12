# Protobuf Source Root

Status: Draft v0.1
Last updated: 2026-05-13
Scope: `proto/`

This directory is the planned source root for vibit Protobuf wire schemas.

The semantic source of truth for business behavior remains under `contracts/` and module manifests. Protobuf files define the client/server wire message shape and must stay aligned with the semantic contract sources.

## Layout

Module Protobuf sources should use:

```text
proto/vibit/<module>/v1/
```

The first planned module path is:

```text
proto/vibit/inventory/v1/
```

Generated Go Protobuf output should go under:

```text
runtime/internal/generated/proto/
```

## Rules

- Run `node tools/vibit check protocol` before creating or changing `.proto` files.
- Do not hand-edit generated Go Protobuf output.
- Keep Protobuf package names, message names, service names, and field names in English.
- Version public wire schemas explicitly.
- Transport adapters should convert Protobuf wire messages into vibit commands and queries; domain modules should not own Protobuf framing.

## Future Tooling

When Protobuf generation starts, root configuration should include:

```text
buf.yaml
buf.gen.yaml
```

Buf linting, formatting, breaking checks, and generation orchestration are accepted by ADR-0013.
