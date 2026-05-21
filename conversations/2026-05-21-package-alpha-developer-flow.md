# Conversation: Package Alpha Developer Flow

Date: 2026-05-21
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0097-package-alpha-developer-flow.md`
Related change: `changes/2026-05-21-package-alpha-developer-flow/`

## Context

The maintainer asked to continue after `W-0188` added the alpha acceptance checklist and was committed and pushed.

The queue identified `M-117/W-0189 Package alpha developer flow` as the next ready work item.

## Maintainer Narrative

The maintainer said:

```text
继续推进
```

Per `docs/workflow.md`, this means advance one next-ready work item.

## Agent Response Summary

The agent added:

- `docs/alpha-developer-flow.md`
- `docs/alpha-developer-flow.zh-CN.md`
- `ADR-0097`
- `runtime.package_alpha_developer_flow` repository check coverage

The packaged flow connects README intake, the runtime runbook, the local request-loop script, status endpoints, alpha acceptance checklist, redaction posture, PostgreSQL manual setup posture, verification commands, and the next contribution entry point.

The slice did not publish `v0.1 alpha`, add release packaging, create release tags or binaries, add hosted deployment, change runtime behavior, change protocol routes, change Protobuf sources or generated output, add migrations, add dependencies, broaden operations/admin behavior, change authentication/session behavior, add product modules, or select direct Nakama/Pitaya API compatibility.

## Decisions

- Accept `ADR-0097`.
- Complete `M-117/W-0189`.
- Add `runtime.package_alpha_developer_flow` as the repository check rule.
- Move the next ready alpha work item to a release publishing decision gate without publishing a release.

## Artifacts

- `docs/alpha-developer-flow.md`
- `docs/alpha-developer-flow.zh-CN.md`
- `decisions/ADR-0097-package-alpha-developer-flow.md`
- `changes/2026-05-21-package-alpha-developer-flow/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Whether to publish `v0.1 alpha` remains deferred to a later explicit decision gate.
- Release packaging remains deferred.
- Public local onboarding protocol route remains deferred.
- Hosted deployments remain deferred.

## Follow-Up

- Advance `W-0190 Release publishing decision gate` as the next bounded contribution step.
- Keep release publishing, release packaging, hosted deployments, runtime behavior changes, protocol changes, Protobuf/generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct Nakama/Pitaya API compatibility, and broad product modules behind later explicit work items.

## Redaction Notes

No raw device credentials, raw access tokens, digest bytes, HMAC input or output, verifier key material, concrete verifier key set ids, database credentials, headers, cookies, query strings, subprotocol values, remote addresses, or GitHub tokens were recorded.
