# ADR-0093: Runtime Runbook Alpha Path Refresh

Status: Accepted
Date: 2026-05-21
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-21-refresh-runtime-runbook-for-alpha/`

Related conversations:

- `conversations/2026-05-21-runtime-runbook-alpha-path-refresh.md`

Related artifacts:

- `docs/runtime-runbook.md`
- `docs/runtime-runbook.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `AGENTS.md`
- `runtime/AGENTS.md`
- `README.md`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0092` proved the local authenticated gameplay alpha path with a focused protocol-level Go E2E test. The existing runtime runbook still described the earlier bootstrap state: memory inventory requests, PostgreSQL startup configuration, and stale notes that public onboarding and presence query were still unavailable.

The `v0.1 alpha` target needs a runbook that describes the actual current path before a minimal example client, local request-loop script, health/config surface, or alpha acceptance checklist can be accurate.

## Decision

Refresh the runtime runbook around the now-proven local alpha path:

```text
local onboarding -> device credential login -> connection binding -> protected inventory -> protected presence -> logout -> post-logout protected request rejection
```

The refreshed runbook:

- distinguishes the default memory runtime path from the PostgreSQL alpha runtime path,
- states that memory startup remains inventory-bootstrap-only,
- states that PostgreSQL startup wires authentication, route protection, connection binding, runtime sessions, logout, and presence query,
- records the verifier key environment variables and key handling requirements,
- documents the focused authenticated gameplay E2E proof command,
- states that local onboarding exists as an application service method and is not yet a public protocol route,
- records redaction rules for raw credentials, raw tokens, digests, HMAC material, verifier keys, concrete key set ids, DSNs, and transport metadata.

This slice also adds `runtime.runtime_runbook_alpha_path_refresh` as repository check coverage for the runbook state and moves the work queue to the next alpha-enabling item: a minimal example client or request-loop script.

## Alternatives Considered

- Leave the runbook stale until a public onboarding tool exists.
- Add an example client or request-loop script before refreshing the runbook.
- Document local onboarding as if it were already a public protocol route.
- Change startup behavior so the runbook can offer a simpler flow.
- Add migrations, protocol routes, generated output, or dependencies as part of the documentation refresh.

## Rationale

Runbook accuracy is now the highest-leverage next step because the alpha path exists as executable proof but is not yet easy for a new developer to understand. The runbook should tell the truth about what is public runtime behavior, what is test-proven application service behavior, and what remains deferred.

Documenting local onboarding as application-service-only avoids overclaiming production signup or public account creation. It also gives the future example client/request-loop work a precise target: make the proven path easier to exercise without changing security boundaries implicitly.

Nakama and Pitaya both imply a developer expectation that authentication, connection lifecycle, gameplay calls, and logout form a coherent local loop. vibit adapts that expectation through its own agent-native documentation and verification boundaries, not through direct API compatibility.

## Agent Reasoning Summary

The maintainer asked to continue. `W-0185` was next ready. Because the task was documentation-focused, the safest slice was to update the runbook and repository checks without touching runtime code, protocol definitions, generated output, migrations, or dependencies.

## Decision Weights

```yaml
decision_weights:
  alpha_developer_usability: high
  documentation_accuracy: high
  security_redaction_clarity: high
  protocol_surface_restraint: high
  runtime_behavior_change: low
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `docs/runtime-runbook.md` and `docs/runtime-runbook.zh-CN.md` now describe the current local alpha path.
- `runtime.runtime_runbook_alpha_path_refresh` becomes the repository check rule for this slice.
- The next alpha work can focus on a minimal example client or request-loop script.
- No runtime behavior, startup semantics, protocol routes, Protobuf sources, generated output, migrations, dependencies, release artifacts, production signup, broad product modules, or direct Nakama/Pitaya API compatibility are added.

## Reversal Conditions

Revisit this decision if local onboarding becomes a public protocol route, if the memory runtime path gains authentication, if startup begins applying migrations automatically, if request-token proof is replaced by session or bound identity for ordinary protected routes, or if direct Nakama/Pitaya API compatibility is explicitly selected.

## Follow-Up

- Add a minimal example client or request-loop script for the proven alpha path.
- Add health/readiness/version/config surfaces.
- Add an alpha acceptance checklist or equivalent check once the developer flow is stable.
