# Conversation: Health Readiness Version Config Surface

Date: 2026-05-21
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0095-health-readiness-version-config-surface.md`
Related change: `changes/2026-05-21-add-health-readiness-version-config-surface/`

## Context

The maintainer asked to continue after `W-0186` added the minimal local alpha request-loop script.

The queue identified `M-115/W-0187 Add health/readiness/version/config surface` as the next ready item.

## Maintainer Narrative

The maintainer said:

```text
继续推进 并进行提交和推送。
```

Per `docs/workflow.md`, this means advance one next-ready work item. The maintainer also explicitly requested committing and pushing the completed state.

## Agent Response Summary

The agent added minimal HTTP status endpoints to the runtime process:

- `/healthz`
- `/readyz`
- `/version`
- `/configz`

The endpoints expose only local troubleshooting posture and avoid secret values. Tests verify status responses and redaction of DSN, verifier key, token audience, and transport-like values.

The slice did not add broad operations/admin behavior, observability dependencies, Protobuf sources, generated output, migrations, release artifacts, broad product modules, or direct Nakama/Pitaya API compatibility.

## Decisions

- Accept `ADR-0095`.
- Complete `M-115/W-0187`.
- Add `runtime.health_readiness_version_config_surface` as the repository check rule.
- Move the next ready alpha work item to an alpha acceptance checklist.

## Artifacts

- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/main_test.go`
- `changes/2026-05-21-add-health-readiness-version-config-surface/`
- `decisions/ADR-0095-health-readiness-version-config-surface.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `docs/runtime-runbook.md`
- `docs/runtime-runbook.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- The alpha acceptance checklist remains deferred.
- Production observability remains deferred.
- Metrics backend integration remains deferred.
- Admin console behavior remains deferred.

## Follow-Up

Advance the alpha acceptance checklist work item.

## Redaction Notes

No raw device credentials, raw access tokens, digest bytes, HMAC input or output, verifier key material, concrete verifier key set ids, database credentials, headers, cookies, query strings, subprotocol values, remote addresses, or GitHub tokens were recorded.
