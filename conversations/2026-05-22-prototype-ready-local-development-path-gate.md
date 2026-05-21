# Conversation: Prototype-Ready Local Development Path Gate

Date: 2026-05-22
Participants: Maintainer, Agent
Related decision: `decisions/ADR-0107-prototype-ready-local-development-path-gate.md`
Related change: `changes/2026-05-22-define-prototype-ready-local-development-path-gate/`

Related artifacts:

- `docs/prototype-ready-local-development-path-gate.md`
- `docs/prototype-ready-local-development-path-gate.zh-CN.md`
- `decisions/ADR-0107-prototype-ready-local-development-path-gate.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `tools/vibit`

## Context

After `W-0198`, the repository had a prototype-ready foundation execution plan. The maintainer asked to continue moving vibit toward a real product and productivity stage.

## Maintainer Narrative

The durable maintainer intent is:

- vibit should move beyond a `0.1` alpha skeleton;
- the project should become a foundation developers can use for game/backend prototypes;
- the next step should reduce local setup, migration, configuration, and example-flow friction;
- the project should remain honest about what is not production-ready.

## Agent Response Summary

The agent completed `W-0199` by defining the prototype-ready local development path gate.

The work added:

- an English and Simplified Chinese gate document;
- `ADR-0107`;
- check coverage for `runtime.prototype_ready_local_development_path_gate`;
- `.arch` continuation updates;
- the next work item `W-0200 Implement prototype-ready local development path package`.

## Decisions

- Accept `ADR-0107`.
- Complete `M-127/W-0199`.
- Select `prototype_ready_local_development_path_package` as the next direction.
- Use `W-0200` to package a repeatable source-first local development path.
- Preserve runtime, protocol, generated-output, migration, dependency, release artifact, hosted deployment, public announcement, paid promotion, operations/admin, authentication/session, broad product, and direct compatibility deferrals.

## Artifacts

- `docs/prototype-ready-local-development-path-gate.md`
- `docs/prototype-ready-local-development-path-gate.zh-CN.md`
- `decisions/ADR-0107-prototype-ready-local-development-path-gate.md`
- `changes/2026-05-22-define-prototype-ready-local-development-path-gate/`
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

- The exact local path package files remain for `W-0200` to implement.
- Whether the later shared-service implementation should start with storage objects or realtime messaging remains a later decision.
- Broader public outreach still requires explicit maintainer authorization.

## Follow-Up

- Continue with `W-0200 Implement prototype-ready local development path package`.
- Keep the first package focused on setup, migrations, configuration, example flow, and verification ergonomics.
- Ask before any runtime behavior, protocol, migration, dependency, release artifact, hosted deployment, announcement, broad product, authentication/session, or direct compatibility expansion.

## Redaction Notes

No raw device credentials, raw access tokens, digest bytes, HMAC input or output, verifier key material, concrete verifier key set ids, database credentials, headers, cookies, query strings, subprotocol values, remote addresses, private local environment file content, or GitHub tokens were recorded.

## Boundary

This work did not implement runtime behavior, change protocol routes, add Protobuf or generated output, add migrations, add dependencies, broaden operations/admin behavior, change authentication/session behavior, add hosted deployment, create release artifacts, execute public announcements, run paid promotion, add broad product modules, or add direct Nakama/Pitaya API compatibility.

## Outcome

`W-0199` is now the durable local development path gate. The next step is to implement the local development path package.
