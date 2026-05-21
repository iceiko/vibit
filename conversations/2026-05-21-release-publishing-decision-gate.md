# Conversation: Release Publishing Decision Gate

Date: 2026-05-21
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0098-release-publishing-decision-gate.md`
Related change: `changes/2026-05-21-define-release-publishing-decision-gate/`

## Context

The maintainer asked to continue after `W-0189` packaged the alpha developer flow and was committed and pushed.

The queue identified `M-118/W-0190 Define release publishing decision gate` as the next ready work item.

## Maintainer Narrative

The maintainer said:

```text
继续推进
```

Per `docs/workflow.md`, this means advance one next-ready work item.

## Agent Response Summary

The agent added:

- `docs/release-publishing-decision-gate.md`
- `docs/release-publishing-decision-gate.zh-CN.md`
- `ADR-0098`
- `runtime.release_publishing_decision_gate` repository check coverage

The gate defines release-publishing prerequisites, release artifact boundaries, verification requirements, stop conditions, redaction expectations, and the next bounded direction.

The slice did not publish `v0.1 alpha`, create release tags, binaries, archives, containers, packages, checksums, hosted deployments, runtime behavior, protocol routes, Protobuf sources or generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct Nakama/Pitaya API compatibility, or broad product modules.

## Decisions

- Accept `ADR-0098`.
- Complete `M-118/W-0190`.
- Add `runtime.release_publishing_decision_gate` as the repository check rule.
- Move the next ready alpha work item to a release execution preparation gate without publishing a release.

## Artifacts

- `docs/release-publishing-decision-gate.md`
- `docs/release-publishing-decision-gate.zh-CN.md`
- `decisions/ADR-0098-release-publishing-decision-gate.md`
- `changes/2026-05-21-define-release-publishing-decision-gate/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Whether and when to execute release publication remains deferred to a later explicit work item.
- Release artifacts remain deferred.
- Hosted deployments remain deferred.
- Public release support posture remains deferred to future preparation work.

## Follow-Up

- Advance `W-0191 Define release execution preparation gate` as the next bounded contribution step.
- Keep release publishing, release artifacts, hosted deployments, runtime behavior changes, protocol changes, Protobuf/generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct Nakama/Pitaya API compatibility, and broad product modules behind later explicit work items.

## Redaction Notes

No raw device credentials, raw access tokens, digest bytes, HMAC input or output, verifier key material, concrete verifier key set ids, database credentials, headers, cookies, query strings, subprotocol values, remote addresses, or GitHub tokens were recorded.
