# Request

## Original Request

```text
继续推进
```

## Clarified Requirement

Advance `W-0185 Refresh runtime runbook for alpha path` from `.arch/work-items.yaml`.

Refresh `docs/runtime-runbook.md` and `docs/runtime-runbook.zh-CN.md` around the now-proven local authenticated gameplay alpha path.

## User-Visible Outcome

A developer or agent can read the runtime runbook and understand:

- the difference between the default memory runtime path and the PostgreSQL alpha runtime path,
- which environment variables are needed for the PostgreSQL path,
- how verifier key material must be handled,
- what the proven local alpha flow is,
- which E2E test proves that flow,
- why local onboarding is not yet a public protocol route,
- which secrets must not be logged or committed.

## Non-Goals

- Do not add runtime behavior.
- Do not change startup configuration semantics.
- Do not add protocol routes.
- Do not change Protobuf sources or generated output.
- Do not add migrations.
- Do not add dependencies.
- Do not publish a release.
- Do not choose production signup, external identity providers, password login, account recovery, account merge, or multi-device linking.
- Do not add chat, social, matchmaking, match runtime, or broad product modules.
- Do not select direct Nakama/Pitaya API compatibility.
- Do not add the future example client or request-loop script in this slice.

## Unknowns

- The concrete local onboarding developer surface remains deferred to a later work item.
- The exact shape of the future example client or request-loop script remains deferred.
- Health/readiness/version/config surfaces and alpha acceptance checks remain deferred.

## Acceptance Criteria

- [x] `docs/runtime-runbook.md` describes the current local alpha path accurately.
- [x] `docs/runtime-runbook.zh-CN.md` is updated with the same meaning.
- [x] The runbook distinguishes in-memory test proof from PostgreSQL-backed runtime operation.
- [x] The runbook records verifier key, raw credential, raw token, digest, and DSN redaction expectations.
- [x] The change does not add runtime behavior, protocol routes, generated output, migrations, dependencies, release artifacts, production identity behavior, broad product modules, or direct Nakama/Pitaya API compatibility.
