# ADR-0094: Minimal Alpha Request Loop Script

Status: Accepted
Date: 2026-05-21
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-21-add-minimal-example-client-or-request-loop-script/`

Related conversations:

- `conversations/2026-05-21-minimal-alpha-request-loop-script.md`

Related artifacts:

- `examples/local-alpha-request-loop.sh`
- `docs/runtime-runbook.md`
- `docs/runtime-runbook.zh-CN.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `README.md`
- `README.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0092` proved the local authenticated gameplay alpha path through a focused Go E2E test. `ADR-0093` refreshed the runtime runbook around that path and left the next developer-experience gap clear: a developer still needed a small local entry point that exercises the proven path without requiring live PostgreSQL, committed secrets, or a hand-built WebSocket client.

The current public protocol still does not expose local onboarding as a public route. A full external client would either need a new onboarding surface or would have to skip the beginning of the alpha path. That would cross the `W-0186` ask-first boundaries.

## Decision

Add `examples/local-alpha-request-loop.sh` as the minimal local alpha request-loop script.

The script invokes the existing focused E2E proof:

```bash
cd runtime && go test ./internal/platform/protocol/protobuf -run TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout -v
```

It prints a short redacted path summary and then lets the existing test runner report pass/fail status. It does not print raw device credentials, raw access tokens, digest material, verifier keys, DSNs, or concrete transport metadata.

This is intentionally a request-loop proof wrapper, not a product SDK, not a public onboarding client, and not a live PostgreSQL process client.

## Alternatives Considered

- Add a full WebSocket client.
- Add public protocol onboarding so an external client can create the first credential.
- Add a PostgreSQL seed tool or CLI onboarding command.
- Add a new Go command under `runtime/cmd`.
- Leave the E2E test as the only documented entry point.

## Rationale

The existing E2E test is the only complete path that includes local onboarding without adding a public onboarding protocol route. Wrapping it in a tiny `examples/` script gives developers a stable, visible entry point while preserving all current security and protocol deferrals.

This keeps the slice inside `W-0186`: it improves alpha usability without changing runtime behavior, startup semantics, protocol surfaces, migrations, generated output, dependencies, release posture, or product scope.

Nakama and Pitaya both set an expectation that developers can quickly run an authenticate-then-play loop. vibit meets that expectation here through a redacted local proof wrapper and keeps direct API compatibility out of scope.

## Agent Reasoning Summary

The maintainer asked to continue. `W-0186` was next ready. The smallest viable implementation was a local script over the existing authenticated gameplay E2E proof because a live external client would require public onboarding or seed behavior not authorized by this work item.

## Decision Weights

```yaml
decision_weights:
  alpha_developer_usability: high
  boundary_restraint: high
  redaction_safety: high
  protocol_surface_restraint: high
  runtime_behavior_change: low
  direct_api_compatibility: low
confidence: high
```

## Consequences

- Developers can run `examples/local-alpha-request-loop.sh` to exercise the proven local alpha path.
- The script is covered by `runtime.minimal_example_client_or_request_loop`.
- The next alpha work can focus on health/readiness/version/config surfaces.
- No runtime behavior, startup semantics, protocol routes, Protobuf sources, generated output, migrations, dependencies, release artifacts, production signup, broad product modules, or direct Nakama/Pitaya API compatibility are added.

## Reversal Conditions

Revisit this decision if local onboarding becomes a public protocol route, if a controlled local seed path is added, if a real SDK/example-client layer is selected, or if the alpha acceptance criteria require a live WebSocket client rather than a request-loop proof wrapper.

## Follow-Up

- Add health/readiness/version/config surfaces.
- Add an alpha acceptance checklist or equivalent check after the developer flow stabilizes.
