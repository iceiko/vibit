# Request

## Original Request

The maintainer asked:

```text
继续推进5小时。
```

## Clarified Requirement

Interpret the maintainer's message as authorization to close the blocked `M-040/W-0112` next-direction confirmation gate and select the agent's conservative recommended direction:

```text
expose_access_token_protocol_carrier_and_route_protection_gate
```

The change must close `M-040`, mark `W-0112` completed with the selected direction, create a bounded next milestone, and create exactly one `next_ready` work item for the selected gate.

## User-Visible Outcome

`node tools/vibit inspect next --json` should no longer show `W-0112` as blocked. The work queue should move to the access-token protocol carrier and route-protection gate milestone with `W-0113` as the next ready work item.

## Non-Goals

- Do not expose protocol proof carriers in code.
- Do not add route-protection implementation.
- Do not add `.proto` files or generated Protobuf output.
- Do not change the existing Protobuf envelope.
- Do not add WebSocket handshake authentication.
- Do not add session persistence.
- Do not wire authentication into process startup.
- Do not change repositories, migrations, dependencies, logout, refresh, cleanup, or token rotation behavior.

## Acceptance Criteria

- [x] `W-0112` is marked completed with `selected_direction: expose_access_token_protocol_carrier_and_route_protection_gate`.
- [x] `M-040` is marked completed.
- [x] The next milestone is created and active.
- [x] `W-0113` is created as the first `next_ready` work item.
- [x] Manifests record the selected direction.
- [x] A conversation log records the maintainer authorization.
- [x] Verification is recorded.
