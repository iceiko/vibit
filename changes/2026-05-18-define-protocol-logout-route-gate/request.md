# Request

## Original Request

The maintainer asked to replan the project route around the Nakama/Pitaya-class product target and continue development.

## Clarified Requirement

After ratifying the product parity roadmap, define the next lifecycle-closure gate for a future client-facing protocol logout route.

## User-Visible Outcome

Future agents can inspect exactly how `runtime.authentication.LogoutAccessToken` should become a WebSocket Protobuf command route without confusing token logout with socket close, session revocation, reconnect, or direct Nakama/Pitaya API compatibility.

## Non-Goals

- Do not add Protobuf logout messages in this change.
- Do not generate Go Protobuf output in this change.
- Do not register a runtime logout route in this change.
- Do not change WebSocket transport, the existing Protobuf envelope, socket close behavior, session revocation behavior, reconnect behavior, session carriers, dependencies, or direct Nakama/Pitaya API compatibility.

## Unknowns

- Whether future logout success should expose `token_record_id` in all production surfaces or only in audit/admin surfaces.
- Whether bound-connection identity should continue after logout until a protected route revalidates.
- Whether a later transport close handoff will notify the client before close or close silently.

## Acceptance Criteria

- [x] The English gate standard exists.
- [x] The Simplified Chinese translation exists.
- [x] ADR-0079 records the decision.
- [x] The work queue marks `W-0169` completed and creates `W-0170` as the next ready implementation slice.
- [x] Repository checks include `runtime.protocol_logout_route_gate`.
