# Conversation: Agent-Native Feature Request Scaffolding Implementation

Date: 2026-05-24

## Context

The maintainer asked to continue toward Nakama and clarified that vibit's product purpose is AI-native development and AI-native testing. `W-0228` defined the feature request scaffolding gate and opened `W-0229`.

## Maintainer Narrative

The maintainer wants users to state backend requirements and have AI carry the work through requirement specification, tests, implementation, verification, and durable repository memory. Nakama is the primary product capability reference. Pitaya remains deferred.

## Agent Response Summary

The agent implemented source-first feature request scaffolding:

- added `changes/_template/feature-request/`;
- added the six required template artifacts;
- added `docs/agent-native-feature-request-scaffolding.md` and Simplified Chinese translation;
- added `tools/vibit scaffold feature`;
- added focused repository check coverage under `runtime.agent_native_feature_request_scaffolding_implementation`;
- accepted `ADR-0137`;
- opened `M-158/W-0230 Pilot scaffolded Nakama feature request intake`.

## Decisions

- `ADR-0137`: implement feature request scaffolding with source-first templates and a narrow deterministic scaffold command.

## Artifacts

- `changes/_template/feature-request/`
- `docs/agent-native-feature-request-scaffolding.md`
- `docs/agent-native-feature-request-scaffolding.zh-CN.md`
- `changes/2026-05-24-implement-agent-native-feature-request-scaffolding/`
- `decisions/ADR-0137-agent-native-feature-request-scaffolding-implementation.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Future users may need richer scaffold options, but the first shape is intentionally deterministic and source-first.

## Follow-Up

- `M-158/W-0230 Pilot scaffolded Nakama feature request intake` should use the new scaffold on a bounded Nakama-style request before broad module expansion.

## Redaction Notes

No raw credentials, raw access tokens, verifier keys, credential or token digests, HMAC inputs or outputs, PostgreSQL DSNs with credentials, GitHub tokens, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or private user data were recorded.
