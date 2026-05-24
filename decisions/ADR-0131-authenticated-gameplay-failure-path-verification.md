# ADR-0131: Authenticated Gameplay Failure Path Verification

Status: Accepted
Date: 2026-05-24
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-24-strengthen-authenticated-gameplay-failure-path-verification/`

Related conversations:

- `conversations/2026-05-24-authenticated-gameplay-failure-path-verification.md`

Related artifacts:

- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
- `docs/agent-native-feature-request-test-workflow.md`
- `docs/nakama-pitaya-product-parity-roadmap.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/contracts.yaml`
- `.arch/conventions.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Context

`ADR-0130` completed the self-presence/status local proof hardening and opened `M-151/W-0223` as the next bounded follow-up. The repository already had unit coverage for several route-authentication failures and a local alpha E2E proof for the authenticated happy path, but the local alpha proof did not yet group the important protected gameplay failure paths under the same real service, route-protection, FrameHandler, and in-memory repository fixture.

Authentication and session behavior are core Nakama-style game backend foundation concerns. Before broader product capability work depends on protected requests, vibit should prove that missing, malformed, invalid, expired, or revoked proof fails closed and remains redacted at the protocol-flow boundary.

Nakama capability family:

```text
identity_authentication_sessions
```

Pitaya remains deferred as a future distributed architecture reference.

## Decision

Strengthen authenticated gameplay failure-path verification using tests only.

The local alpha proof adds `TestAuthenticatedGameplayFailurePathsLocalAlphaFlow` and covers:

```text
protected inventory without authenticated wrapper
protected inventory with malformed authenticated wrapper
protected inventory with malformed access-token text
protected inventory with unknown well-formed access-token text
protected inventory with expired access token
protected inventory with revoked access token after logout
protected presence without authenticated wrapper
error redaction for raw credential and access-token text
```

Open:

```text
M-152/W-0224 Select next Nakama prototype-ready capability after authenticated failure-path proof
```

as the next-ready work item.

This decision does not add runtime behavior beyond tests, protocol routes, Protobuf source, generated output, migrations, dependencies, startup wiring, persistence, login behavior changes, logout behavior changes, token refresh, cleanup jobs, session validation changes, chat, groups, matchmaking, match runtime, distributed runtime, hosted deployments, release artifacts, public announcements, paid promotion, or direct Nakama/Pitaya API compatibility.

## Alternatives Considered

- Add a new protocol route or proof carrier to represent failure testing explicitly.
- Change route policy or authentication service behavior while adding the tests.
- Add only lower-level unit tests and skip the real local alpha FrameHandler path.
- Jump directly to chat, groups, matchmaking, or match runtime after the presence proof.
- Treat Pitaya distributed routing as part of the near-term failure-path work.

## Rationale

The selected work item asked for verification hardening, not feature expansion. The existing local alpha fixture can exercise the actual authentication service, token validator, route protector, protocol frame handler, inventory route, presence route, logout route, and in-memory repositories. That boundary gives stronger product confidence than another isolated unit test while still avoiding protocol or runtime scope expansion.

The next step should be a direction-selection slice because the immediate authenticated failure-path gap is closed. Selecting the next Nakama prototype-ready capability keeps the project aligned with the maintainer's Nakama-first direction while avoiding an unbounded jump into broad modules.

## Agent Reasoning Summary

The maintainer asked to keep advancing toward Nakama and emphasized AI-native development and testing. The active work item required strengthening failure-path proof. The codebase showed that production behavior already existed, so the correct small step was to add an E2E test around existing behavior and record durable memory.

## Decision Weights

```yaml
decision_weights:
  nakama_product_alignment: high
  alpha_foundation_confidence: high
  failure_path_testability: high
  redaction_confidence: high
  bounded_scope: high
  reuse_existing_protocol_surface: high
  direct_api_compatibility: low
confidence: high
```

## Consequences

- `M-151/W-0223` is completed.
- `M-152/W-0224` is next-ready.
- Authenticated local alpha failure-path proof now covers missing, malformed, invalid, expired, revoked, route-protection, and redaction behavior.
- The proof uses existing runtime and Protobuf surfaces.
- Broad authentication/session changes, protocol changes, token refresh, cleanup jobs, chat, groups, matchmaking, match runtime, persistence, distributed runtime, SDKs, operations, and direct compatibility remain deferred.
- `runtime.authenticated_gameplay_failure_path_verification` checks the verification records and next-ready state.

## Reversal Conditions

Revisit this decision if:

- a future external client fixture replaces in-process Protobuf frame tests as the preferred alpha proof boundary;
- token/session policy changes alter the expected public error mapping;
- route policy moves protected gameplay routes from request-token proof to a different identity requirement;
- alpha users report a more urgent Nakama-style blocker than the next selected prototype-ready capability.

## Follow-Up

- Complete `W-0224`: select the next Nakama prototype-ready capability after authenticated failure-path proof.
- Keep new protocol routes, Protobuf, generated output, migrations, dependencies, token refresh, cleanup jobs, broad modules, distributed runtime, and direct compatibility behind later bounded work items.
