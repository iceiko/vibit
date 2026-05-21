# ADR-0095: Health Readiness Version Config Surface

Status: Accepted
Date: 2026-05-21
Decision Makers: Maintainer, Agent
Related changes:

- `changes/2026-05-21-add-health-readiness-version-config-surface/`

Related conversations:

- `conversations/2026-05-21-health-readiness-version-config-surface.md`

Related artifacts:

- `runtime/cmd/vibit-server/main.go`
- `runtime/cmd/vibit-server/main_test.go`
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

The local alpha path now has a runbook and a minimal request-loop script, but a developer running the process still lacks a small status surface for answering basic troubleshooting questions:

- Is the process alive?
- Is it ready?
- What runtime version is this?
- Which runtime store posture is active?
- Where is the WebSocket endpoint?
- Are secrets being redacted?

The work item explicitly forbids broad operations/admin behavior, observability dependencies, secret disclosure, generated output, migrations, release publishing, broad product modules, and direct Nakama/Pitaya API compatibility.

## Decision

Add minimal HTTP status endpoints to `runtime/cmd/vibit-server`:

- `/healthz`
- `/readyz`
- `/version`
- `/configz`

The status surface uses the existing `net/http` mux and standard library JSON encoding. `/configz` reports only redacted posture: runtime store, WebSocket path, local alpha request-loop script path, PostgreSQL configured boolean, authentication configuration required boolean, and `secrets_redacted: true`.

The surface must not expose raw device credentials, raw access tokens, verifier key values, DSNs, digests, headers, cookies, query strings, subprotocol values, remote addresses, or concrete transport metadata.

## Alternatives Considered

- Add a full operations/admin API.
- Add metrics endpoints and observability dependencies.
- Add Protobuf status messages.
- Add a CLI config inspection command.
- Leave status information only in runbook text.

## Rationale

HTTP status endpoints are the smallest useful local troubleshooting surface because the runtime already owns an HTTP mux for `/v1/ws`. They help developers confirm process state without changing gameplay protocol, authentication behavior, persistence behavior, generated output, migrations, dependencies, or release posture.

Nakama and Pitaya both set expectations that a local server process should be inspectable. vibit meets that expectation with a small redacted local surface rather than copying their API shapes or adding broad operations behavior.

## Agent Reasoning Summary

The maintainer asked to continue and also requested commit and push. `W-0187` was next ready. The implementation stayed in `runtime/cmd/vibit-server` because startup and process HTTP mux wiring already live there. Tests verify both endpoint shape and redaction.

## Decision Weights

```yaml
decision_weights:
  alpha_troubleshooting_value: high
  redaction_safety: high
  implementation_size: small
  operations_scope_restraint: high
  dependency_addition: none
  direct_api_compatibility: none
confidence: high
```

## Consequences

- Developers can query `/healthz`, `/readyz`, `/version`, and `/configz` on the runtime HTTP server.
- The config endpoint reports posture but not secret values.
- `runtime.health_readiness_version_config_surface` becomes the repository check rule.
- The next alpha work can focus on an acceptance checklist.

## Reversal Conditions

Revisit this decision if alpha packaging selects a different operations surface, if versioning moves to generated build metadata, if a CLI-only troubleshooting model is chosen, or if a production operations standard supersedes the local status endpoints.

## Follow-Up

- Add an alpha acceptance checklist or equivalent repository check.
