# Conversation: Scaffolded Nakama Feature Request Intake Pilot

Date: 2026-05-24

## Context

`W-0229` implemented source-first feature request scaffolding through `tools/vibit scaffold feature` and opened `M-158/W-0230 Pilot scaffolded Nakama feature request intake`.

## Maintainer Narrative

The maintainer asked to continue toward the Nakama target and previously clarified that Pitaya should be deferred for now. The maintainer also clarified that the product purpose is AI-native development and AI-native testing: a user states a requirement, and AI should help produce the spec, acceptance criteria, tests, implementation, verification, and durable memory.

## Agent Response Summary

The agent completed W-0230 by using the new scaffold on one bounded Nakama-style intake:

```text
changes/2026-05-24-pilot-scaffolded-nakama-feature-request-intake/
```

The scaffolded intake selected a concrete Nakama-style request:

```text
player_friendship_relationship_lifecycle
```

The selected capability family is:

```text
friends_groups_and_parties
```

The selected follow-up is:

```text
W-0231 Define friends relationship lifecycle gate
```

No runtime behavior, protocol route, Protobuf source, generated output, migration, dependency, persistence, startup wiring, SDK, hosted deployment, release artifact, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility was added.

## Decisions

- `ADR-0138`: complete scaffolded Nakama feature request intake pilot, select `friends_groups_and_parties`, and open `W-0231 Define friends relationship lifecycle gate`.
- Check rule: `runtime.scaffolded_nakama_feature_request_intake_pilot`.

## Artifacts

- `changes/2026-05-24-pilot-scaffolded-nakama-feature-request-intake/`
- `decisions/ADR-0138-scaffolded-nakama-feature-request-intake-pilot.md`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/contracts.yaml`
- `.arch/conventions.yaml`
- `tools/vibit`
- `rules/check-rules.json`

## Open Questions

- W-0231 must decide the first semantic boundary, ownership, invariants, errors, events, permissions, data posture, and protocol posture for friendship lifecycle before implementation.

## Follow-Up

Complete W-0231 as a gate only. It should define semantic contracts, acceptance criteria, test plan, ownership, data posture, protocol posture, generated output posture, redaction, verification, and stop conditions before any friendship implementation.

## Redaction Notes

No raw credentials, raw access tokens, verifier keys, credential or token digests, HMAC inputs or outputs, PostgreSQL DSNs with credentials, GitHub tokens, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, or private user data beyond the explicit feature request text were recorded.
