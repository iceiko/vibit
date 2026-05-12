# Protobuf Source Root

Status: Draft v0.1
Last updated: 2026-05-13
Scope: `proto/`

This directory is the planned source root for vibit Protobuf wire schemas.

The semantic source of truth for business behavior remains under `contracts/` and module manifests. Protobuf files define the client/server wire message shape and must stay aligned with the semantic contract sources.

The game protocol envelope standard is defined by `docs/game-protocol.md`, `.arch/protocol.yaml`, and `ADR-0015`. Read those artifacts before creating `.proto` files. The first protocol model is a WebSocket-framed Protobuf envelope with explicit `kind`, `module`, and `name` routing fields, session metadata, target metadata, server-authoritative message rules, and error mapping.

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
- Keep envelope and module payload schemas aligned with `.arch/protocol.yaml`.
- Do not hand-edit generated Go Protobuf output.
- Keep Protobuf package names, message names, service names, and field names in English.
- Version public wire schemas explicitly.
- Transport adapters should convert Protobuf wire messages into vibit commands and queries; domain modules should not own Protobuf framing.
- Do not implement room state sync, matchmaking, allocation, reconnect replay, presence, streams, realtime input, or state patches in Protobuf files before their modules and standards exist.

## Future Tooling

When Protobuf generation starts, root configuration should include:

```text
buf.yaml
buf.gen.yaml
```

Buf linting, formatting, breaking checks, and generation orchestration are accepted by ADR-0013.
