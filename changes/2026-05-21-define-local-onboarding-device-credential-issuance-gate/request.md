# Request

## Original Request

The maintainer asked:

```text
继续推进
```

## Clarified Requirement

Advance one `next_ready` work item from `.arch/work-items.yaml`. The queue identifies `W-0181 Define local onboarding device credential issuance gate` as the next ready work item.

This change must define the local onboarding/device credential issuance gate before implementation.

## User-Visible Outcome

Maintainers and future agents can read the gate, ADR, manifests, and check rule to understand exactly how the next implementation may create a local player account and issue the first device credential.

## Non-Goals

- Implement onboarding or credential issuance behavior.
- Generate or display raw credentials.
- Create player accounts through a new runtime flow.
- Write credential records through a new runtime flow.
- Add public protocol routes, Protobuf sources, or generated output.
- Add migrations or dependencies.
- Change credential schema or repository interfaces.
- Add production signup, external identity providers, password login, account recovery, account merge, or multi-device linking.
- Publish a release.
- Add direct Nakama/Pitaya API compatibility.

## Unknowns

- The exact future local tooling or CLI surface remains deferred.
- The concrete player id, event id, and credential record id generation strategy remains deferred.
- Production signup behavior remains deferred.

## Acceptance Criteria

- [x] The gate is documented in English and Simplified Chinese.
- [x] ADR-0089 records the decision.
- [x] The future implementation owner, allowed files, repository handoff, helper composition, one-time raw credential presentation, redaction, tests, and deferrals are explicit.
- [x] The work queue marks `W-0181` completed and creates a bounded `W-0182` implementation item.
- [x] Repository checks cover the gate.
