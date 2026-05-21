# Conversation: Next Alpha Direction Local Onboarding Device Credential Issuance Gate

Date: 2026-05-21
Participants: Maintainer, Agent
Related change:

- `changes/2026-05-21-select-local-onboarding-device-credential-issuance-gate/`

Related decision:

- `decisions/ADR-0088-next-alpha-direction-local-onboarding-device-credential-issuance-gate.md`

## Context

`W-0178` completed the protected presence protocol query. The queue then moved to `M-108 Next Alpha Direction Selection`, and `node tools/vibit inspect next` reported no `next_ready` work item.

The short-term target is `v0.1 alpha`. The alpha goal document recommends local onboarding/device credential issuance after the presence protocol query.

## Maintainer Narrative

The maintainer said:

```text
继续推进
```

Per `docs/workflow.md`, this means advance one continuation step unless blocked. Because no `next_ready` work item existed, the active step was to complete the direction selection milestone.

## Agent Response Summary

The agent selected `define_local_onboarding_device_credential_issuance_gate` as the next alpha-enabling direction.

The selected next work is gate-only because onboarding/device credential issuance is security-sensitive. It may eventually create player accounts, generate raw device credentials, compute credential verifier digests, and store digest-only credential records. The gate must define redaction, one-time presentation, failure behavior, repository ordering, and deferrals before implementation.

## Decisions

- Complete `M-108/W-0180`.
- Accept `ADR-0088`.
- Create `M-109/W-0181` as the next ready gate-only work item.
- Keep runtime onboarding behavior, protocol messages, generated output, migrations, dependencies, release publishing, and direct Nakama/Pitaya API compatibility deferred.

## Artifacts

- `changes/2026-05-21-select-local-onboarding-device-credential-issuance-gate/`
- `decisions/ADR-0088-next-alpha-direction-local-onboarding-device-credential-issuance-gate.md`
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

## Open Questions

- Whether the first local onboarding surface should be CLI, documented local script, protocol route, or another runtime-owned local-only surface remains for `W-0181`.
- Whether player account creation and credential issuance are one service method or two coordinated local operations remains for `W-0181`.
- Whether the first implementation should support only PostgreSQL runtime store remains for `W-0181`.

## Follow-Up

Run the next continuation against `W-0181` to define the gate before any local onboarding or device credential issuance code is added.

## Redaction Notes

No raw access tokens, device credentials, generated secrets, digest bytes, HMAC input bytes, verifier keys, database secrets, player private data, headers, cookies, query strings, subprotocol values, remote addresses, or GitHub tokens are recorded in this conversation log.
