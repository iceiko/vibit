# Conversation: Dependency Adoption Standard

Date: 2026-05-12
Participants: Maintainer, Agent
Related changes:

- `changes/2026-05-12-add-dependency-adoption-standard/`

Related artifacts:

- `.arch/dependencies.yaml`
- `docs/dependency-adoption.md`
- `docs/dependency-adoption.zh-CN.md`
- `docs/_templates/dependency-adoption.md`
- `tools/vibit`

## Context

After PostgreSQL and S3-compatible object storage were ratified as the persistence direction, the next likely branch points were concrete dependency choices: WebSocket library, Protobuf tooling, PostgreSQL driver, migration tooling, S3 SDK, MinIO deployment, test framework, and observability stack.

## Maintainer Narrative

The maintainer asked:

```text
按照你的理解继续，有任何可能产生分支的地方，与我讨论确认。
```

This reinforced the rule that the agent may continue necessary preparation, but should not silently choose a direction when multiple plausible branches exist.

## Agent Response Summary

The agent continued with a conservative preparation step: create a dependency adoption standard, template, and `.arch/dependencies.yaml` registry before choosing any specific dependency.

The change records dependency slots and required evaluation criteria without selecting concrete packages such as a PostgreSQL driver, migration tool, WebSocket library, Protobuf toolchain, or S3 SDK.

## Decisions

- Add a dependency adoption standard governed by `ADR-0010`.
- Add `.arch/dependencies.yaml` as the machine-readable registry for foundational dependency decision slots.
- Keep all concrete dependency choices unselected until explicit adoption records are created.
- Require future agents to discuss and confirm branch-producing dependency choices before marking a slot `accepted`.

## Artifacts

- Added `docs/dependency-adoption.md`.
- Added `docs/dependency-adoption.zh-CN.md`.
- Added `docs/_templates/dependency-adoption.md`.
- Added `.arch/dependencies.yaml`.
- Updated architecture README, repository README, AGENTS guides, and `tools/vibit check architecture`.
- Added `changes/2026-05-12-add-dependency-adoption-standard/`.

## Open Questions

- Which WebSocket library should be adopted?
- Which Protobuf toolchain and generation layout should be adopted?
- Which PostgreSQL driver should be adopted?
- Which migration tool or migration convention should be adopted?
- Whether first Go tests should use only the standard library or adopt a test helper library.

## Follow-Up

- Discuss concrete dependency candidates with the maintainer before accepting any dependency slot.
- Use the dependency adoption template for the next dependency decision.
- Add import-boundary checks after Go runtime code exists.

## Redaction Notes

No secrets, tokens, account identifiers, or private data were included in this conversation log.
