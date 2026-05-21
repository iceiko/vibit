# Conversation: Runtime Runbook Alpha Path Refresh

Date: 2026-05-21
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0093-runtime-runbook-alpha-path-refresh.md`
Related change: `changes/2026-05-21-refresh-runtime-runbook-for-alpha/`

## Context

The maintainer asked to continue after `W-0184` proved the authenticated gameplay E2E path.

The queue identified `M-113/W-0185 Refresh runtime runbook for alpha path` as the next ready item.

## Maintainer Narrative

The maintainer said:

```text
继续推进
```

Per `docs/workflow.md`, this means advance one next-ready work item unless blocked.

## Agent Response Summary

The agent refreshed the runtime runbook and its Simplified Chinese translation around the current local alpha path:

```text
local onboarding -> login -> connection binding -> protected inventory grant/read -> protected presence query -> logout -> post-logout protected request rejection
```

The runbook now distinguishes the default memory inventory bootstrap path from the PostgreSQL alpha runtime path, records verifier key environment requirements, points to the focused E2E proof test, states that local onboarding is not yet a public protocol route, and records redaction rules for raw credentials, raw tokens, digests, verifier keys, DSNs, and transport metadata.

No runtime behavior, startup semantics, protocol routes, Protobuf sources, generated output, migrations, dependencies, release artifacts, production signup, broad product modules, or direct Nakama/Pitaya API compatibility were added.

## Decisions

- Accept `ADR-0093`.
- Complete `M-113/W-0185`.
- Add `runtime.runtime_runbook_alpha_path_refresh` as the repository check rule.
- Move the next ready alpha work item to a minimal example client or request-loop script.

## Artifacts

- `docs/runtime-runbook.md`
- `docs/runtime-runbook.zh-CN.md`
- `changes/2026-05-21-refresh-runtime-runbook-for-alpha/`
- `decisions/ADR-0093-runtime-runbook-alpha-path-refresh.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `docs/v0.1-alpha-goal.md`
- `docs/v0.1-alpha-goal.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- The minimal example client or local request-loop script remains deferred.
- Health/readiness/version/config surfaces remain deferred.
- The alpha acceptance checklist remains deferred.
- Public protocol onboarding remains deferred.

## Follow-Up

Advance the minimal example client or request-loop script work item.

## Redaction Notes

No raw device credentials, raw access tokens, digest bytes, HMAC input or output, verifier key material, concrete verifier key set ids, database credentials, headers, cookies, query strings, subprotocol values, remote addresses, or GitHub tokens were recorded.
