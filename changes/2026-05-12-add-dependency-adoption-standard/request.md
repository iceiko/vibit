# Request

Date: 2026-05-12
Change ID: `add-dependency-adoption-standard`
Type: standard

## Maintainer Request

The maintainer asked the agent to continue according to its judgment, while discussing and confirming any place where a branch in direction could occur:

```text
按照你的理解继续，有任何可能产生分支的地方，与我讨论确认。
```

## Clarified Requirement

Continue runtime readiness work without selecting concrete dependencies yet.

Create a standard and registry that make future dependency choices explicit before they are ratified.

## User-Visible Outcome

Future agents should have a clear dependency adoption process before choosing:

- WebSocket library.
- Protobuf tooling.
- PostgreSQL driver.
- Migration tooling.
- S3 SDK.
- MinIO deployment dependency.
- Test framework beyond Go standard tooling.
- Observability stack.

## Non-Goals

- Do not choose a concrete dependency.
- Do not add Go runtime code.
- Do not add package manifests.
- Do not accept any dependency slot.
- Do not supersede Go, WebSocket, Protobuf, PostgreSQL, or S3-compatible object-storage decisions.

## Acceptance Criteria

- [x] Dependency adoption standard exists in English and Simplified Chinese.
- [x] Dependency adoption template exists.
- [x] `.arch/dependencies.yaml` records foundational dependency slots and statuses.
- [x] `.arch/conventions.yaml` references the dependency adoption standard, template, and registry.
- [x] README and AGENTS files reference the new dependency adoption artifacts.
- [x] Architecture check covers the new required artifacts.
