# Impact

This is a direction-selection and queue-management change only.

It completes `M-108/W-0180` and opens `M-109/W-0181` as the next gate-only alpha slice:

```text
define_local_onboarding_device_credential_issuance_gate
```

## Affected Scope

- Work queue state.
- Runtime architecture manifests.
- Alpha goal and continuation documentation.
- Agent handoff context.

## No Runtime Behavior Added

This change does not add:

- local onboarding execution,
- player account creation behavior,
- device credential issuance behavior,
- credential repository writes from onboarding,
- raw credential presentation,
- Protobuf messages,
- generated output,
- migrations,
- routes,
- startup wiring,
- CLI commands,
- HTTP endpoints,
- WebSocket behavior,
- dependencies,
- release artifacts,
- or direct Nakama/Pitaya API compatibility.

## Security Posture

The selected next step is gate-only because onboarding/device credential issuance touches credential material, verifier digests, one-time raw secret presentation, player account creation, and repository mutation ordering. The implementation must be defined before code is added.
