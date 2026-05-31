# Request

## Original Request

Continue the next ready repository work item.

## Clarified Requirement

Advance `M-168/W-0240 Define friends relationship protocol route gate`.

Define only the future protocol route gate for exposing the friends relationship runtime behavior service from `ADR-0147`.

## User-Visible Outcome

Future protocol implementation work has a documented route family, message-shape posture, protected-route identity handoff, public error mapping, service handoff, generated-output posture, and stop conditions before any `.proto`, generated output, route handlers, or startup wiring are added.

## Non-Goals

- No protocol route implementation.
- No Protobuf source.
- No generated output.
- No dependency additions.
- No migration changes.
- No startup wiring.
- No authentication/session behavior changes.
- No event/audit table.
- No delivery guarantees.
- No stream subscriptions.
- No chat, groups, parties, broadcast fanout, matchmaking, match runtime, SDK, hosted deployment, release artifact, distributed runtime, or direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- The gate document and Simplified Chinese translation exist.
- `ADR-0148` accepts the gate.
- The future route family maps to `runtime/internal/app/friends.Service` methods.
- The gate requires `request_token_required`, authenticated wrapper handoff, and validated request identity as the only actor proof.
- The gate records candidate `vibit.friends.v1` Protobuf shapes without adding them.
- The gate maps `FRIENDSHIP_*` public service errors for future protocol use.
- Repository checks, manifests, module docs, README/alpha docs, and work queue state identify W-0241 as the next bounded implementation slice.
