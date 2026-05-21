# Conversation: Minimal Alpha Request Loop Script

Date: 2026-05-21
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0094-minimal-alpha-request-loop-script.md`
Related change: `changes/2026-05-21-add-minimal-example-client-or-request-loop-script/`

## Context

The maintainer asked to continue after `W-0185` refreshed the runtime runbook around the proven local alpha path.

The queue identified `M-114/W-0186 Add minimal example client or request-loop script` as the next ready item.

## Maintainer Narrative

The maintainer said:

```text
继续推进
```

Per `docs/workflow.md`, this means advance one next-ready work item unless blocked.

## Agent Response Summary

The agent added `examples/local-alpha-request-loop.sh` as the smallest developer-facing local alpha request-loop entry point. The script runs the existing focused authenticated gameplay E2E proof and prints only a redacted path summary plus test output.

The slice did not add runtime behavior, startup semantics, public protocol onboarding, Protobuf sources, generated output, migrations, dependencies, release artifacts, production signup, broad product modules, or direct Nakama/Pitaya API compatibility.

## Decisions

- Accept `ADR-0094`.
- Complete `M-114/W-0186`.
- Add `runtime.minimal_example_client_or_request_loop` as the repository check rule.
- Move the next ready alpha work item to health/readiness/version/config surfaces.

## Artifacts

- `examples/local-alpha-request-loop.sh`
- `changes/2026-05-21-add-minimal-example-client-or-request-loop-script/`
- `decisions/ADR-0094-minimal-alpha-request-loop-script.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `docs/runtime-runbook.md`
- `docs/runtime-runbook.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `runtime/AGENTS.md`
- `runtime/AGENTS.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Health/readiness/version/config surfaces remain deferred.
- The alpha acceptance checklist remains deferred.
- Public protocol onboarding remains deferred.
- A live external WebSocket example client remains deferred until the onboarding or seed path is explicitly selected.

## Follow-Up

Advance the health/readiness/version/config surface work item.

## Redaction Notes

No raw device credentials, raw access tokens, digest bytes, HMAC input or output, verifier key material, concrete verifier key set ids, database credentials, headers, cookies, query strings, subprotocol values, remote addresses, or GitHub tokens were recorded.
