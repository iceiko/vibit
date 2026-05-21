# Conversation: First Alpha User Discovery Loop

Date: 2026-05-22
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0104-first-alpha-user-discovery-loop.md`
Related change: `changes/2026-05-21-define-first-alpha-user-discovery-loop/`
Related changes:

- `changes/2026-05-21-define-first-alpha-user-discovery-loop/`

Related artifacts:

- `docs/first-alpha-user-discovery-loop.md`
- `docs/first-alpha-user-discovery-loop.zh-CN.md`
- `decisions/ADR-0104-first-alpha-user-discovery-loop.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `tools/vibit`

## Context

The maintainer previously authorized `W-0195` as a `go` for the source-first `v0.1.0-alpha.1` release and asked that the README become more attractive so vibit can find users. The release has been pushed and published as a Git tag, GitHub release record, release notes, and GitHub source archive only.

The continuation queue then exposed `W-0196 Define first alpha user discovery loop`.

## Maintainer Narrative

The durable maintainer intent is:

- continue after the source-first alpha release;
- find real users;
- keep the README and release surface attractive and honest;
- do not expand release artifacts, hosted deployment, runtime behavior, protocol routes, generated output, migrations, dependencies, operations/admin behavior, authentication/session behavior, broad product modules, or direct Nakama/Pitaya API compatibility without explicit authorization.

## Agent Response Summary

The agent completed `W-0196` by defining the first alpha user discovery loop.

The loop records:

- target developer segments;
- allowed and deferred outreach surfaces;
- feedback capture fields;
- review questions;
- success signals;
- stop conditions;
- the next ready work item, `W-0197 Prepare first alpha feedback intake surfaces`.

## Decisions

- Accept `ADR-0104`.
- Complete `M-124/W-0196`.
- Define the first alpha user discovery loop as a learning loop, not a growth campaign.
- Record target developer segments for Go game/backend, Nakama/Pitaya evaluators, AI-agent backend engineers, open-source contributors, and source-first prototype builders.
- Record allowed and deferred outreach surfaces before any broad announcement.
- Record feedback fields, review questions, success signals, and stop conditions.
- Open `M-125/W-0197 Prepare first alpha feedback intake surfaces` as the next ready direction.
- Preserve public announcement, paid promotion, hosted deployment, release artifact, runtime, protocol, generated-output, migration, dependency, operations/admin, authentication/session, broad product module, and direct compatibility deferrals.

## Artifacts

- `docs/first-alpha-user-discovery-loop.md`
- `docs/first-alpha-user-discovery-loop.zh-CN.md`
- `decisions/ADR-0104-first-alpha-user-discovery-loop.md`
- `changes/2026-05-21-define-first-alpha-user-discovery-loop/`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `README.md`
- `README.zh-CN.md`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- Which feedback intake surface should be prepared first: issue templates, discussion guidance, labels, or a smaller repository-owned note?
- Whether and when the maintainer will authorize public announcements beyond the GitHub release record remains undecided.

## Follow-Up

- Continue with `W-0197 Prepare first alpha feedback intake surfaces`.
- Keep intake guidance redaction-safe.
- Ask for explicit maintainer authorization before broad outreach, paid promotion, hosted deployment, or additional release artifacts.

## Redaction Notes

No raw device credentials, raw access tokens, digest bytes, HMAC input or output, verifier key material, concrete verifier key set ids, database credentials, headers, cookies, query strings, subprotocol values, remote addresses, private local environment file content, or GitHub tokens were recorded.

## Boundary

This work did not execute public announcements beyond the GitHub release record, paid promotion, hosted deployments, additional release artifacts, runtime behavior changes, protocol changes, Protobuf or generated output changes, migrations, dependency changes, broad operations/admin behavior, authentication/session changes, broad product module expansion, or direct Nakama/Pitaya API compatibility.

## Outcome

`W-0196` is now the durable planning record for finding early developers. The next step is to prepare repository-owned feedback intake surfaces before any broader outreach.
