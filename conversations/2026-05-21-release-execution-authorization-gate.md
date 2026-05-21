# Conversation: Release Execution Authorization Gate

Date: 2026-05-21
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0100-release-execution-authorization-gate.md`
Related change: `changes/2026-05-21-define-release-execution-authorization-gate/`

## Context

The maintainer asked to continue after `W-0191` defined the release execution preparation gate and was committed and pushed.

The queue identified `M-120/W-0192 Define release execution authorization gate` as the next ready work item.

## Maintainer Narrative

The maintainer said:

```text
继续推进
```

Per `docs/workflow.md`, this means advance one next-ready work item unless blocked or confirmation is required.

## Agent Response Summary

The agent added:

- `docs/release-execution-authorization-gate.md`
- `docs/release-execution-authorization-gate.zh-CN.md`
- `ADR-0100`
- `runtime.release_execution_authorization_gate` repository check coverage

The gate defines final go/no-go criteria, required verification state, release identifier review, artifact authorization boundaries, maintainer approval requirements, authorization outcome, stop conditions, redaction expectations, and the next bounded maintainer decision point.

The slice did not publish `v0.1 alpha`, select or create release identifiers or tags, create release binaries, archives, containers, packages, checksums, provenance files, hosted deployments, runtime behavior, protocol routes, Protobuf sources or generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct Nakama/Pitaya API compatibility, or broad product modules.

## Decisions

- Accept `ADR-0100`.
- Complete `M-120/W-0192`.
- Add `runtime.release_execution_authorization_gate` as the repository check rule.
- Create `M-121/W-0193 Confirm release execution maintainer decision` as a blocked maintainer decision gate.

## Artifacts

- `docs/release-execution-authorization-gate.md`
- `docs/release-execution-authorization-gate.zh-CN.md`
- `decisions/ADR-0100-release-execution-authorization-gate.md`
- `changes/2026-05-21-define-release-execution-authorization-gate/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Whether release execution is go or no-go remains deferred to explicit maintainer authorization.
- The final release identifier remains deferred.
- Release artifact creation remains deferred.
- Hosted deployments remain deferred.

## Follow-Up

- Stop at `W-0193 Confirm release execution maintainer decision` until the maintainer explicitly chooses go/no-go and authorizes any release identifier, artifact, or publication boundary.
- Keep release publishing, release artifacts, hosted deployments, runtime behavior changes, protocol changes, Protobuf/generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct Nakama/Pitaya API compatibility, and broad product modules behind later explicit work items.

## Redaction Notes

No raw device credentials, raw access tokens, digest bytes, HMAC input or output, verifier key material, concrete verifier key set ids, database credentials, headers, cookies, query strings, subprotocol values, remote addresses, or GitHub tokens were recorded.
