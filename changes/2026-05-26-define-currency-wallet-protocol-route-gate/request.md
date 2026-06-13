# Request

## Original Request

Continue up to twenty next-ready repository work items.

## Clarified Requirement

Advance `M-229/W-0301 Define currency wallet protocol route gate`.

Define only the future protocol route gate for exposing the currency wallet runtime behavior service from `ADR-0208`.

## User-Visible Outcome

Future protocol implementation work has a documented route family, message-shape posture, protected-route identity handoff, public error mapping, service handoff, generated-output posture, local proof expectations, and stop conditions before any `.proto`, generated output, route handlers, or startup wiring are added.

## Non-Goals

- No protocol route implementation.
- No Protobuf source.
- No generated output.
- No dependency additions.
- No migration changes.
- No startup wiring.
- No authentication/session behavior changes.
- No reward integration.
- No inventory integration.
- No purchase behavior.
- No catalog table.
- No event/audit table.
- No payment behavior.
- No reservation, settlement, refund, transfer, operations/admin, SDK, hosted deployment, release artifact, distributed runtime, or direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- The gate document and Simplified Chinese translation exist.
- `ADR-0209` accepts the gate.
- The future route family maps to `runtime/internal/app/currency.Service` methods.
- The gate requires `request_token_required`, authenticated wrapper handoff, and validated request identity as the only owner proof.
- The gate rejects metadata-only identity and client-supplied owner/wallet/actor proof.
- The gate records candidate `vibit.currency.v1` Protobuf shapes without adding them.
- The gate maps `CURRENCY_WALLET_*` public service errors for future protocol use.
- Repository checks, manifests, module docs, README/alpha docs, and work queue state identify W-0302 as the next bounded implementation slice.
