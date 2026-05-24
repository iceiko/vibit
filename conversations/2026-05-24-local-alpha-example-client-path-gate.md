# Conversation: Local Alpha Example Client Path Gate

Date: 2026-05-24

## Context

The maintainer asked the agent to continue toward the Nakama target for many steps. `M-152/W-0224` selected `client_sdks_examples_and_developer_experience` as the next bounded Nakama prototype-ready capability family and opened `M-153/W-0225` as a gate.

## Maintainer Narrative

The maintainer wants Nakama as the primary product reference. Pitaya is deferred. The product purpose remains AI-native development and AI-native testing: user requirements should become specs, acceptance criteria, tests, implementation, verification, and durable memory through AI assistance.

## Agent Response Summary

The agent defined the local alpha example client path gate. The gate selects a source-first repository-local example path instead of a public SDK or live external client surface.

The future implementation candidates are:

```text
examples/local-alpha-client/README.md
examples/local-alpha-example-client.sh
runtime/internal/platform/protocol/protobuf/authenticated_gameplay_e2e_test.go
```

The gate opens:

```text
M-154/W-0226 Implement local alpha example client path
```

as the next-ready bounded follow-up.

## Decisions

- Accept `ADR-0133`.
- Accept `docs/local-alpha-example-client-path-gate.md`.
- Register `runtime.local_alpha_example_client_path_gate`.
- Complete `M-153/W-0225`.
- Open `M-154/W-0226 Implement local alpha example client path` as next-ready.
- Keep Pitaya deferred as a future distributed architecture reference.
- Keep direct Nakama/Pitaya API compatibility out of scope.

## Artifacts

- `docs/local-alpha-example-client-path-gate.md`
- `docs/local-alpha-example-client-path-gate.zh-CN.md`
- `changes/2026-05-24-define-local-alpha-example-client-path-gate/`
- `decisions/ADR-0133-local-alpha-example-client-path-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Should `W-0226` only add docs and a wrapper script, or should it also add a smaller named local alpha example-flow test next to the existing Protobuf E2E tests?
- Should a later public client package boundary move generated Protobuf output out of `runtime/internal/`?

## Follow-Up

- Complete `W-0226`: implement the local alpha example client path inside the gate.

## Redaction Notes

No raw access-token, device credential, verifier digest, lookup digest, verifier key, DSN, local secret, or GitHub token value is recorded here.
