# Conversation: Alpha Acceptance Checklist

Date: 2026-05-21
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0096-alpha-acceptance-checklist.md`
Related change: `changes/2026-05-21-add-alpha-acceptance-checklist/`

## Context

The maintainer asked to continue after `W-0187` added the health/readiness/version/config troubleshooting surface and was committed and pushed.

The queue identified `M-116/W-0188 Add alpha acceptance checklist` as the next ready work item.

## Maintainer Narrative

The maintainer said:

```text
继续推进
```

Per `docs/workflow.md`, this means advance one next-ready work item.

## Agent Response Summary

The agent added:

- `docs/alpha-acceptance-checklist.md`
- `docs/alpha-acceptance-checklist.zh-CN.md`
- `ADR-0096`
- `runtime.alpha_acceptance_checklist` repository check coverage

The checklist records ready, manual, deferred, and blocked acceptance states for the local alpha developer flow. It covers repository intake, prerequisites, configuration, migration posture, runtime surfaces, authenticated gameplay, verification commands, redaction, and contribution entry points.

The slice did not publish `v0.1 alpha`, add release packaging, change runtime behavior, change protocol routes, change Protobuf sources or generated output, add migrations, add dependencies, broaden operations/admin behavior, change authentication/session behavior, add product modules, or select direct Nakama/Pitaya API compatibility.

## Decisions

- Accept `ADR-0096`.
- Complete `M-116/W-0188`.
- Add `runtime.alpha_acceptance_checklist` as the repository check rule.
- Move the next ready alpha work item to packaging the alpha developer flow without release publishing.

## Artifacts

- `docs/alpha-acceptance-checklist.md`
- `docs/alpha-acceptance-checklist.zh-CN.md`
- `decisions/ADR-0096-alpha-acceptance-checklist.md`
- `changes/2026-05-21-add-alpha-acceptance-checklist/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Release publishing remains deferred.
- Release packaging remains deferred.
- Public local onboarding protocol route remains deferred.
- Broad operations/admin behavior remains deferred.

## Follow-Up

- Advance `W-0189 Package alpha developer flow` as the next bounded contribution step.
- Keep release publishing, release packaging, hosted deployments, runtime behavior changes, protocol changes, Protobuf/generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior changes, direct Nakama/Pitaya API compatibility, and broad product modules behind later explicit work items.

## Redaction Notes

No raw device credentials, raw access tokens, digest bytes, HMAC input or output, verifier key material, concrete verifier key set ids, database credentials, headers, cookies, query strings, subprotocol values, remote addresses, or GitHub tokens were recorded.
