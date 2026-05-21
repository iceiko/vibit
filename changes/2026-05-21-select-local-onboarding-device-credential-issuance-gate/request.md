# Request

Advance the active `M-108 Next Alpha Direction Selection` milestone after the protected presence protocol query is complete.

The maintainer asked to continue. `node tools/vibit inspect next` reports no `next_ready` item, while `.arch/work-items.yaml` shows `M-108` is active and asks agents to select the next bounded alpha-enabling work item.

## Selected Direction

```text
define_local_onboarding_device_credential_issuance_gate
```

## Rationale

`docs/v0.1-alpha-goal.md` lists first local onboarding/device credential issuance as the next preferred alpha step after the presence protocol query. That work is security-sensitive because it creates player account state, generates a raw device credential for one-time client visibility, stores digest-only credential verifier records, and defines an initial local developer onboarding path.

Per `ADR-0082`, security-critical work requires a separate gate before implementation. This change therefore selects and queues a gate-only work item. It does not implement onboarding or credential issuance.

## Non-Goals

- Do not add runtime onboarding behavior.
- Do not generate or display device credential material.
- Do not create player accounts from a new onboarding flow.
- Do not write credential verifier records from a new onboarding flow.
- Do not add Protobuf messages, generated output, routes, migrations, dependencies, startup wiring, CLI behavior, HTTP endpoints, WebSocket behavior, or release artifacts.
- Do not publish `v0.1 alpha`.
- Do not add direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- [x] `define_local_onboarding_device_credential_issuance_gate` is selected.
- [x] `M-108/W-0180` records the direction selection.
- [x] `M-109/W-0181` becomes the next ready gate-only work item.
- [x] Security-sensitive implementation remains deferred behind the new gate.
