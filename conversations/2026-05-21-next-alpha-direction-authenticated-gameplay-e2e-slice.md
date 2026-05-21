# Conversation: Next Alpha Direction Authenticated Gameplay E2E Slice

Date: 2026-05-21
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0091-next-alpha-direction-authenticated-gameplay-e2e-slice.md`
Related change: `changes/2026-05-21-select-next-alpha-direction-after-local-onboarding/`

## Context

The maintainer asked to continue after `W-0182` completed the local-only onboarding/device credential issuance implementation.

The queue identified `M-111/W-0183 Select next alpha direction after local onboarding` as the next ready item. The alpha goal document recommends authenticated gameplay E2E after local onboarding.

## Maintainer Narrative

The maintainer said:

```text
继续推进
```

Per `docs/workflow.md`, this means advance one next-ready work item unless blocked.

## Agent Response Summary

The agent selected `define_authenticated_gameplay_e2e_slice` as the next alpha-enabling direction after local onboarding.

This selection creates `M-112/W-0184` as the next ready bounded functional slice. The intended flow is local onboarding -> login -> connection binding -> protected inventory -> presence query -> logout.

The direction-selection step did not implement authenticated gameplay E2E, add protocol routes, change Protobuf sources, change generated output, add migrations, add dependencies, publish a release, add production signup, add external identity providers, add password login, add account recovery, add multi-device linking, add broad product modules, or add direct Nakama/Pitaya API compatibility.

## Decisions

- Accept `ADR-0091`.
- Complete `M-111/W-0183`.
- Create `M-112/W-0184 Define authenticated gameplay E2E slice` as the next ready work item.
- Keep runtime runbook refresh, minimal example client/request-loop script, health/readiness/version/config surface, and alpha acceptance checklist queued after the E2E direction.

## Artifacts

- `changes/2026-05-21-select-next-alpha-direction-after-local-onboarding/`
- `decisions/ADR-0091-next-alpha-direction-authenticated-gameplay-e2e-slice.md`
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

- Whether `W-0184` should be primarily a Go E2E test, a local request-loop script, documentation, or a small combination remains for the W-0184 change spec.
- Whether W-0184 should require PostgreSQL-only verification or keep memory-store coverage remains for the W-0184 change spec.
- Whether any protocol surface is missing for E2E proof remains for W-0184 analysis.

## Follow-Up

Advance `M-112/W-0184` to define and prove the authenticated gameplay E2E slice.

## Redaction Notes

No raw device credentials, raw access tokens, digest bytes, HMAC input or output, verifier key material, database credentials, player private data, headers, cookies, query strings, subprotocol values, remote addresses, or GitHub tokens were recorded.
