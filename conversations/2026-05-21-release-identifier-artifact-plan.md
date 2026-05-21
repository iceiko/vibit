# Conversation: Release Identifier And Artifact Plan

Date: 2026-05-21
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0102-release-identifier-artifact-plan.md`
Related change: `changes/2026-05-21-define-release-identifier-artifact-plan/`

## Context

The queue was next-ready at `M-122/W-0194 Define release identifier and artifact plan` after `W-0193` recorded approval to continue release execution planning.

The plan needed to define a proposed release identifier, check for conflicts without creating public release state, and define future artifact boundaries before any release-producing command could be considered.

## Maintainer Narrative

The maintainer said:

```text
继续推进
```

## Agent Response Summary

The agent interpreted the request as authorization to complete the next ready planning item only. The agent did not treat it as authorization to publish `v0.1 alpha`, create tags, push tags, create release artifacts, create GitHub release records, or run release publication commands.

## Decisions

- Accept `ADR-0102`.
- Complete `M-122/W-0194`.
- Record `v0.1.0-alpha.1` as the proposed release identifier.
- Record no observed local tag, remote origin tag, or GitHub release record conflict for `v0.1.0-alpha.1`.
- Define the proposed future source-first release surface as Git tag, GitHub release record, repository-derived release notes, and hosting-platform source archive.
- Defer binaries, packages, containers, checksums, provenance, hosted deployments, package registry publication, generated SDK packages, install scripts, and public announcements beyond a release record.
- Create blocked `M-123/W-0195 Confirm release execution final authorization`.

## Artifacts

- `docs/release-identifier-artifact-plan.md`
- `docs/release-identifier-artifact-plan.zh-CN.md`
- `decisions/ADR-0102-release-identifier-artifact-plan.md`
- `changes/2026-05-21-define-release-identifier-artifact-plan/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Final release execution approval remains ungranted.
- Exact final release commit remains unselected.
- Tag annotation or signing posture remains unapproved.
- GitHub release creation command or UI/API workflow remains unapproved.

## Follow-Up

- Use `W-0195` to record final maintainer go/no-go authorization before any release-producing command runs.
- Repeat tag and release-record conflict checks immediately before any execution step.

## Redaction Notes

No release was published, no release identifier was selected for execution, no tag was created or pushed, no GitHub release record was created, no release artifact was created, and no runtime behavior, protocol route, Protobuf source, generated output, migration, dependency, operations/admin behavior, authentication/session behavior, broad product module, hosted deployment, or direct Nakama/Pitaya API compatibility was added.

No raw device credentials, raw access tokens, digest bytes, HMAC input or output, verifier key material, concrete verifier key set ids, database credentials, headers, cookies, query strings, subprotocol values, remote addresses, or GitHub tokens were recorded.
