# Conversation: Authenticated Gameplay Failure Path Verification

Date: 2026-05-24

## Context

The maintainer asked the agent to continue for many steps toward the Nakama target. The current next-ready work item was `M-151/W-0223 Strengthen authenticated gameplay failure-path verification`, opened by `ADR-0130`.

The selected Nakama capability family was `identity_authentication_sessions`.

## Maintainer Narrative

The maintainer clarified that the project should target Nakama first and that AI-native development plus AI-native testing are the product purpose. User requirements should become specs, tests, implementation, verification, and durable memory through the AI workflow.

## Agent Response Summary

The agent read the work queue, workflow standard, Nakama roadmap, route protection code, authentication service behavior, FrameHandler behavior, and authenticated local alpha E2E fixture. The agent determined that production behavior already existed and that the correct bounded step was to strengthen tests:

- protected inventory without authenticated wrapper now fails closed in the local alpha E2E proof;
- protected inventory with malformed authenticated wrapper now fails closed;
- malformed access-token text maps to the public malformed-token class;
- unknown well-formed token text maps to the public invalid-token class;
- expired token records map to the public invalid-token class;
- revoked tokens after logout map to the public invalid-token class;
- protected presence without authenticated wrapper now fails closed;
- error envelopes are checked for raw access-token and one-time device credential redaction.

No production runtime behavior, protocol route, Protobuf source, generated output, migration, dependency, persistence, startup wiring, token refresh, cleanup job, authentication/session behavior change, or direct compatibility scope was added.

## Decisions

- Accept `ADR-0131`.
- Register `runtime.authenticated_gameplay_failure_path_verification`.
- Complete `M-151/W-0223`.
- Open `M-152/W-0224 Select next Nakama prototype-ready capability after authenticated failure-path proof` as next-ready.
- Keep Pitaya deferred and keep direct Nakama/Pitaya API compatibility out of scope.

## Artifacts

- `changes/2026-05-24-strengthen-authenticated-gameplay-failure-path-verification/`
- `decisions/ADR-0131-authenticated-gameplay-failure-path-verification.md`
- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/contracts.yaml`
- `.arch/conventions.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Which Nakama prototype-ready capability should follow authenticated gameplay failure-path hardening?
- Should the next slice select a concrete feature family such as notifications, friends/groups, leaderboards, or operations inspection, or should it strengthen developer ergonomics first?

## Follow-Up

- Complete `W-0224`: select the next Nakama prototype-ready capability after authenticated failure-path proof.

## Redaction Notes

No raw access-token, device credential, verifier digest, lookup digest, verifier key, local secret, or GitHub token value is recorded here.
