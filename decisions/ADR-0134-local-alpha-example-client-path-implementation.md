# ADR-0134: Local Alpha Example Client Path Implementation

Status: Accepted
Date: 2026-05-24
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-24-implement-local-alpha-example-client-path/`

Related conversations:

- `conversations/2026-05-24-local-alpha-example-client-path-implementation.md`

Related artifacts:

- `examples/local-alpha-client/README.md`
- `examples/local-alpha-client/README.zh-CN.md`
- `examples/local-alpha-example-client.sh`
- `examples/local-alpha-request-loop.sh`
- `examples/README.md`
- `examples/README.zh-CN.md`
- `docs/local-alpha-example-client-path-gate.md`
- `decisions/ADR-0133-local-alpha-example-client-path-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0133` selected a source-first repository-local local alpha example client path. The repository already proves the alpha loop through focused Go E2E tests, but developers and AI agents need a clearer example entrypoint under `examples/` that demonstrates the current capability loop without pretending vibit has a public SDK or generated client package.

The current constraints still matter: generated Protobuf Go output remains internal to the runtime, local onboarding is not a public client route, and the alpha example must not disclose credentials, access tokens, digests, verifier keys, DSNs, or transport metadata.

## Decision

Implement the first local alpha example client path with:

```text
examples/local-alpha-client/README.md
examples/local-alpha-client/README.zh-CN.md
examples/local-alpha-example-client.sh
```

Keep:

```text
examples/local-alpha-request-loop.sh
```

as a compatibility wrapper around the new script.

The script runs existing focused runtime proofs:

```text
TestAuthenticatedGameplayE2EUsesExistingOnboardingLoginBindingInventoryPresenceAndLogout
TestStorageObjectsProtocolRouteLocalAlphaFlow
TestPresenceStatusLocalAlphaFlowReportsOfflineAfterCloseAndInvalidation
TestAuthenticatedGameplayFailurePathsLocalAlphaFlow
```

Open `M-155/W-0227 Select next Nakama prototype-ready capability after local alpha example client path` as the next-ready follow-up.

This decision does not publish an SDK, generate client libraries, add dependencies, add runtime behavior, add public onboarding routes, add protocol routes, add Protobuf source, change generated output, add migrations, add persistence, change startup wiring, change authentication/session behavior, add hosted deployments, create release artifacts, add Pitaya-style distributed architecture, or add direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Build a live WebSocket client outside `runtime/internal`.
- Publish a first SDK or generated client library.
- Add a public local onboarding route before examples.
- Only update the old `local-alpha-request-loop.sh` without adding a readable example path.
- Jump directly to chat, groups, matchmaking, match runtime, operations/admin, or distributed runtime.

## Rationale

The smallest useful implementation is a readable, redacted, source-first example path that reuses existing proofs. It improves developer experience immediately while keeping the public client boundary honest.

This matches Nakama's product pressure that a backend framework should show a client-facing evaluation path, but vibit keeps its own route names, payloads, token posture, and no-compatibility stance.

## Agent Reasoning Summary

The maintainer asked to keep advancing toward Nakama. The next-ready work item was W-0226, which ADR-0133 bounded to examples and existing runtime proof reuse. The safe implementation is docs plus shell entrypoint, with no runtime or protocol expansion.

## Decision Weights

```yaml
decision_weights:
  nakama_developer_experience_alignment: high
  source_alpha_honesty: high
  implementation_boundedness: high
  redaction_safety: high
  runtime_behavior_risk: none
  direct_api_compatibility: low
confidence: high
```

## Consequences

- The first local alpha example client path exists under `examples/`.
- The old request-loop script remains available as a wrapper.
- `runtime.local_alpha_example_client_path_implementation` becomes the repository check rule for this slice.
- `M-154/W-0226` is completed.
- `M-155/W-0227` becomes next-ready.
- SDK publication, generated client libraries, public onboarding routes, protocol changes, runtime changes, dependencies, hosted demos, release artifacts, distributed runtime, and direct compatibility remain deferred.

## Reversal Conditions

Revisit this decision if:

- a public client package boundary is ratified;
- generated client output moves to a public package;
- local onboarding becomes a public protocol route;
- maintainers explicitly authorize SDK publication or hosted examples;
- alpha feedback shows operations/admin inspection is a stronger immediate blocker.

## Follow-Up

- Complete `W-0227`: select the next Nakama prototype-ready capability after the local alpha example path implementation.
- Keep SDKs, generated clients, hosted demos, protocol changes, runtime changes, dependencies, distributed runtime, and direct compatibility behind later bounded work items.

