# ADR-0100: Release Execution Authorization Gate

Status: Accepted
Date: 2026-05-21
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-21-define-release-execution-authorization-gate/`

Related conversations:

- `conversations/2026-05-21-release-execution-authorization-gate.md`

Related artifacts:

- `docs/release-execution-authorization-gate.md`
- `docs/release-execution-authorization-gate.zh-CN.md`
- `docs/release-execution-preparation-gate.md`
- `docs/release-publishing-decision-gate.md`
- `docs/alpha-developer-flow.md`
- `docs/alpha-acceptance-checklist.md`
- `docs/runtime-runbook.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0099` defined the release execution preparation gate without authorizing release execution or artifact creation. The next bounded step is to define the authorization criteria that must be satisfied before maintainers can make a later go/no-go release execution decision.

`W-0192` was selected to define the release execution authorization gate without publishing a release, selecting a release identifier, creating release tags, or creating release artifacts.

## Decision

Add `docs/release-execution-authorization-gate.md` and `docs/release-execution-authorization-gate.zh-CN.md`.

The gate defines:

- final go/no-go criteria,
- required verification state,
- release identifier review requirements,
- artifact authorization boundaries,
- maintainer approval requirements,
- authorization outcome,
- stop conditions,
- redaction expectations,
- and the next bounded maintainer decision point.

The repository check rule for this gate is `runtime.release_execution_authorization_gate`.

The gate explicitly keeps:

```yaml
release_declared: false
release_publishing_authorized_by_this_gate: false
release_execution_authorized_by_this_gate: false
release_packaging_authorized_by_this_gate: false
release_artifacts_created_by_this_gate: false
hosted_deployment_authorized_by_this_gate: false
release_identifier_selected: false
release_tag_created: false
release_binary_created: false
release_archive_created: false
release_container_created: false
release_package_created: false
release_checksum_created: false
release_provenance_created: false
```

The next bounded work item is `W-0193 Confirm release execution maintainer decision`, and it is intentionally blocked until the maintainer explicitly makes a go/no-go decision and authorizes any release identifier, artifact, or publication boundary.

## Alternatives Considered

- Execute release publication immediately.
- Select the final release identifier inside the authorization gate.
- Create release tags, binaries, archives, packages, containers, checksums, provenance files, or hosted deployments in the authorization slice.
- Treat existing readiness and preparation records as implicit release permission.
- Move directly to broad product features.

## Rationale

Authorization criteria are useful only if they are separated from both execution and maintainer approval. The project now has a durable checklist for a future go/no-go decision while keeping artifact creation and publication behind a later explicit approval record.

Nakama and Pitaya set expectations that server frameworks distinguish release readiness, release decisions, artifacts, and deployment posture. vibit should keep that discipline without copying their APIs, package model, deployment model, operations surfaces, or compatibility promises.

## Agent Reasoning Summary

The maintainer asked to continue. `W-0192` was next ready. Because its ask-first boundary includes publishing, release identifiers, tags, artifacts, and hosted deployments, the implementation stayed documentation- and tooling-focused and ended at a blocked maintainer decision gate.

## Decision Weights

```yaml
decision_weights:
  release_scope_restraint: high
  authorization_clarity: high
  maintainer_approval_clarity: high
  artifact_boundary_clarity: high
  runtime_behavior_change: none
  dependency_addition: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- The repository has an explicit release execution authorization gate.
- Repository checks guard the authorization boundary and its deferrals.
- The next release-path work is blocked on an explicit maintainer go/no-go decision.
- The repository remains pre-alpha and no release identifier, release tag, release artifact, hosted deployment, or published release record is created by this decision.

## Reversal Conditions

Revisit this decision if maintainers choose a different release process, if authorization and approval need to be merged into one record, or if release execution needs separate identifier, artifact, publication, and deployment gates.

## Follow-Up

- Wait for explicit maintainer go/no-go authorization before advancing `W-0193 Confirm release execution maintainer decision`.
- Keep release publication, release tags, binaries, archives, containers, packages, checksums, provenance files, hosted deployments, runtime changes, protocol changes, generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct compatibility, and broad product expansion behind later explicit work items.
