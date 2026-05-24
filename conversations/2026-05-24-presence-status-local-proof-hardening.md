# Conversation: Presence Status Local Proof Hardening

Date: 2026-05-24

## Context

The maintainer asked the agent to continue for many steps toward the Nakama target. The current next-ready work item was `M-150/W-0222 Harden presence status local proof through close and offline cases`, opened by `ADR-0129`.

The selected Nakama capability family was `presence_status_and_notifications`.

## Maintainer Narrative

The maintainer clarified that the project should target Nakama first and that AI-native development plus AI-native testing are the product purpose. User requirements should become specs, tests, implementation, verification, and durable memory through the AI workflow.

## Agent Response Summary

The agent read the work queue, workflow standard, Nakama roadmap, presence service, connection registry, close policy, and authenticated gameplay E2E fixture. The agent determined that runtime behavior already existed and that the correct bounded step was to harden tests:

- application presence service tests now prove offline after transport close and policy invalidation;
- authenticated local alpha Protobuf flow tests now prove online, close/offline, and invalidation/offline using existing protocol surfaces;
- durable memory and checks record `W-0222` completion and open `W-0223`.

No production runtime behavior, protocol route, Protobuf source, generated output, migration, dependency, persistence, startup wiring, subscription, broadcast fanout, or direct compatibility scope was added.

## Decisions

- Accept `ADR-0130`.
- Register `runtime.presence_status_local_proof_hardening`.
- Complete `M-150/W-0222`.
- Open `M-151/W-0223 Strengthen authenticated gameplay failure-path verification` as next-ready.
- Keep Pitaya deferred and keep direct Nakama/Pitaya API compatibility out of scope.

## Artifacts

- `changes/2026-05-24-harden-presence-status-local-proof/`
- `decisions/ADR-0130-presence-status-local-proof-hardening.md`
- `runtime/internal/app/presence/presence_test.go`
- `runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/contracts.yaml`
- `.arch/conventions.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Should future presence/status work add subscriptions or broadcast fanout, or should authenticated gameplay failure-path verification come first?
- What external client fixture, if any, should eventually complement the in-process local alpha request-loop proof?

## Follow-Up

- Complete `W-0223`: strengthen authenticated gameplay failure-path verification.

## Redaction Notes

No raw access-token, device credential, verifier digest, lookup digest, verifier key, local secret, or GitHub token value is recorded here.
