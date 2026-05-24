# Conversation: Agent-Native Feature Request Scaffolding Gate

Date: 2026-05-24

## Context

The maintainer redirected the project to use Nakama as the primary reference and defer Pitaya. The maintainer clarified that the product purpose is AI-native development and AI-native testing: a user requirement should be carried by AI through specification, acceptance criteria, tests, implementation, verification, and durable memory.

`ADR-0135` selected `agent_native_requirement_test_implementation_workflow` as the next Nakama-first prototype-ready capability family and opened `M-156/W-0228 Define agent-native feature request scaffolding gate`.

## Maintainer Narrative

The maintainer wants vibit to progress toward Nakama-class capability coverage, but the near-term differentiator is the AI-native workflow. Users should not have to manually design the whole change process after stating a requirement. The repository should help agents scaffold the right artifacts and tests before implementation starts.

## Agent Response Summary

The agent defined a gate-only boundary for future feature request scaffolding. The gate selects source-first repository-local change artifact scaffolding, records required future artifacts, documents redaction and stop conditions, and opens `M-157/W-0229 Implement agent-native feature request scaffolding`.

No scaffolding implementation, templates, scaffold command, runtime behavior, protocol route, Protobuf source, generated output, migration, dependency, SDK, hosted surface, distributed runtime, or direct compatibility was added.

## Decisions

- Accepted `ADR-0136`.
- Registered `runtime.agent_native_feature_request_scaffolding_gate`.
- Completed `W-0228`.
- Opened `W-0229 Implement agent-native feature request scaffolding`.
- Kept Nakama as the primary product reference.
- Kept Pitaya deferred as a future distributed architecture reference.

## Artifacts

- `docs/agent-native-feature-request-scaffolding-gate.md`
- `docs/agent-native-feature-request-scaffolding-gate.zh-CN.md`
- `decisions/ADR-0136-agent-native-feature-request-scaffolding-gate.md`
- `changes/2026-05-24-define-agent-native-feature-request-scaffolding-gate/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- The exact command-line UX for `tools/vibit scaffold feature` remains for W-0229.
- Whether the scaffold should accept only a change id or also optional title/requirement flags remains for W-0229.
- Whether future scaffolding should include optional module-specific sections remains deferred until a concrete feature request needs it.

## Follow-Up

- Implement `W-0229` within docs, templates, tooling, and checks only.
- Keep runtime behavior, protocol routes, Protobuf source, generated output, migrations, dependencies, SDK publication, hosted demos, distributed runtime, and direct compatibility deferred unless later explicitly authorized.

## Redaction Notes

No local secret values were recorded. Raw credentials, access tokens, verifier keys, lookup digests, verifier digests, HMAC inputs or outputs, PostgreSQL DSNs with credentials, GitHub tokens, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, and concrete transport metadata remain redacted.

