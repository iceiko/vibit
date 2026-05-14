# Request

## Original Request

Continue advancing the project without stopping for non-decision details.

## Clarified Requirement

After ratifying runtime session validation semantic contracts, make `node tools/vibit inspect contract` able to inspect runtime-owned contract families, starting with `runtime:ValidateSession`.

## User-Visible Outcome

Agents can inspect runtime session contracts through the same machine-readable command they already use for domain module contracts.

## Non-Goals

- Do not add new public runtime behavior.
- Do not add authentication, token, credential, persistence, Protobuf, WebSocket handshake, or player runtime implementation.
- Do not change the inspect output schema.
- Do not add a major dependency.

## Acceptance Criteria

- `node tools/vibit inspect contract --module runtime --type command --id ValidateSession` returns JSON.
- Runtime contract inspection reports registry, source, ownership, and consistency fields.
- Existing domain contract inspection still works.
- Repository checks pass.
