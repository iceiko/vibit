# Impact

## Affected Modules

- `runtime`
- `authentication`

## Module Ownership Impact

The future verifier digest comparison helper remains application-owned under `runtime/internal/app/authentication`.

The storage-neutral `authentication` module remains repository-boundary only. It must not compare verifier digests or decide authentication outcomes.

## Public Contract Impact

No public command, query, event, error, permission, Protobuf, or WebSocket contract changes are made.

## Event Impact

No event contracts or runtime event publication behavior are changed.

## Permission Impact

No permissions are added, removed, or changed.

## Data And Migration Impact

No schema or migration changes are made.

## Runtime Impact

No Go runtime comparison code is added. The change only defines future comparison helper boundaries, expected file paths, constant-time primitive posture, redaction rules, test requirements, and deferrals.

## Documentation Impact

Adds:

- `docs/verifier-digest-comparison-helper-gate.md`
- `docs/verifier-digest-comparison-helper-gate.zh-CN.md`
- `decisions/ADR-0049-verifier-digest-comparison-helper-gate.md`
- `conversations/2026-05-15-verifier-digest-comparison-helper-gate.md`

Updates manifests, AGENTS guides, rule catalog, and CLI checks so future agents can discover and verify this gate.

## Compatibility Risks

No runtime compatibility risk. The main risk is over-constraining future implementation names. The selected names are intentionally narrow and can be changed by a later ADR if necessary.
