# Conversation: First Alpha Feedback Intake Surfaces

Date: 2026-05-22
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0105-first-alpha-feedback-intake-and-product-maturity-milestones.md`
Related change: `changes/2026-05-22-prepare-first-alpha-feedback-intake-surfaces/`

Related artifacts:

- `docs/first-alpha-feedback-intake-surfaces.md`
- `docs/first-alpha-feedback-intake-surfaces.zh-CN.md`
- `docs/product-maturity-milestones.md`
- `docs/product-maturity-milestones.zh-CN.md`
- `.github/ISSUE_TEMPLATE/alpha-feedback.yml`
- `decisions/ADR-0105-first-alpha-feedback-intake-and-product-maturity-milestones.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `tools/vibit`

## Context

The maintainer asked what maturity stage vibit's server has reached. The agent answered that `v0.1.0-alpha.1` is a real source-first alpha with an authenticated gameplay loop, but not yet an industrial or production server foundation. The agent described a stage path: source-first alpha, prototype-ready foundation, single-node production-candidate foundation, and Nakama/Pitaya-class product.

The maintainer then asked to set those stages as project milestones, preserve them in the repository, and keep pushing toward the real product and productivity stage.

## Maintainer Narrative

The durable maintainer intent is:

- do not stop at a `0.1` alpha skeleton;
- advance toward a real product stage;
- advance toward a framework that can become actual production leverage for game/backend developers;
- keep the current alpha honest while making the next maturity stages explicit;
- prepare a feedback intake surface so early users can shape the next useful work.

## Agent Response Summary

The agent completed `W-0197` by preparing a repository-owned alpha feedback intake surface and recording product maturity milestones.

The work added:

- a GitHub issue form for alpha feedback;
- an intake standard with redaction, triage labels, maturity buckets, and stop conditions;
- a product maturity milestone standard;
- `ADR-0105`;
- change records and manifest/check updates.

## Decisions

- Accept `ADR-0105`.
- Complete `M-125/W-0197`.
- Use `.github/ISSUE_TEMPLATE/alpha-feedback.yml` as the first feedback intake surface.
- Define product maturity stages:
  - source-first alpha, reached;
  - prototype-ready foundation, next product stage;
  - single-node production-candidate foundation, planned;
  - Nakama/Pitaya-class product, long-term target.
- Set `W-0198 Define prototype-ready foundation execution plan` as the next ready direction.
- Preserve public announcement, paid promotion, hosted deployment, release artifact, runtime, protocol, generated-output, migration, dependency, operations/admin, authentication/session, broad product module, and direct compatibility deferrals.

## Artifacts

- `.github/ISSUE_TEMPLATE/alpha-feedback.yml`
- `docs/first-alpha-feedback-intake-surfaces.md`
- `docs/first-alpha-feedback-intake-surfaces.zh-CN.md`
- `docs/product-maturity-milestones.md`
- `docs/product-maturity-milestones.zh-CN.md`
- `decisions/ADR-0105-first-alpha-feedback-intake-and-product-maturity-milestones.md`
- `changes/2026-05-22-prepare-first-alpha-feedback-intake-surfaces/`
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

- Which prototype-ready capability should lead after the execution plan: setup friction, example client/app, storage objects, realtime messaging/server push, concurrency/failure verification, or operations inspection?
- Whether and when the maintainer will authorize broader public outreach remains undecided.
- Whether GitHub Discussions should be added later remains undecided.

## Follow-Up

- Continue with `W-0198 Define prototype-ready foundation execution plan`.
- Use early user feedback when available, but do not wait indefinitely before defining the first prototype-ready plan.
- Keep intake guidance redaction-safe.
- Ask for explicit maintainer authorization before broad outreach, paid promotion, hosted deployment, additional release artifacts, runtime scope changes, protocol changes, dependencies, broad product modules, or direct compatibility claims.

## Redaction Notes

No raw device credentials, raw access tokens, digest bytes, HMAC input or output, verifier key material, concrete verifier key set ids, database credentials, headers, cookies, query strings, subprotocol values, remote addresses, private local environment file content, or GitHub tokens were recorded.

## Boundary

This work did not execute public announcements beyond the GitHub release record, paid promotion, hosted deployments, additional release artifacts, runtime behavior changes, protocol changes, Protobuf or generated output changes, migrations, dependency changes, broad operations/admin behavior, authentication/session changes, broad product module expansion, or direct Nakama/Pitaya API compatibility.

## Outcome

`W-0197` is now the durable feedback intake and product maturity milestone record. The next step is to define the prototype-ready foundation execution plan.
