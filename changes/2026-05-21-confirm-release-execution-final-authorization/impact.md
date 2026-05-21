# Impact

## Summary

This change records final release execution authorization for the source-first `v0.1.0-alpha.1` alpha, refreshes README as an external user entry point, and prepares release notes for the GitHub release record.

## Public Surface

Authorized public release surface:

- Git tag `v0.1.0-alpha.1`;
- GitHub release record `v0.1.0-alpha.1`;
- GitHub source archive generated from the tag;
- release notes from repository-owned facts.

## Documentation

Adds:

- `docs/release-execution-final-authorization.md`
- `docs/release-execution-final-authorization.zh-CN.md`
- `decisions/ADR-0103-release-execution-final-authorization.md`
- `conversations/2026-05-21-release-execution-final-authorization.md`
- `docs/releases/v0.1.0-alpha.1.md`
- `docs/releases/v0.1.0-alpha.1.zh-CN.md`

Updates README and manifest pointers so external developers see the value proposition, try path, current alpha limits, and continuation state first.

## Runtime, Protocol, Data, Dependencies

No runtime behavior is changed.

No protocol routes are added or changed.

No Protobuf source or generated output is changed.

No migrations are added.

No dependencies are added.

No authentication/session behavior is changed.

## Artifact Deferrals

The change does not authorize binaries, packages, containers, checksums, provenance or signing artifacts, hosted deployments, install scripts, registry publication, SDK packages, or public announcements beyond the GitHub release record.

## Compatibility

No API, event, data, generated-output, or dependency compatibility break is introduced.

Release state changes from pre-release planning to authorized source-first alpha release execution after verification.
