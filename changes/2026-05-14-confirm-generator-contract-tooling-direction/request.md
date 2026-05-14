# Request

## Original Request

The maintainer accepted the recommendation to continue with generator, contract tooling, and architecture check hardening before adding more runtime features, and asked to continue directly without unnecessary confirmation.

## Clarified Requirement

Close the next-direction gate by selecting `improve_generators_and_contract_tooling_before_more_runtime_features`, create an active tooling milestone, and preserve all deferred authentication, token, credential, persistence, protocol handshake, runtime handler, and production permission decisions.

## User-Visible Outcome

`node tools/vibit inspect next --json` should show an active tooling hardening milestone instead of a blocked direction gate.

## Non-Goals

- Do not implement authentication.
- Do not choose token behavior.
- Do not add credential storage.
- Do not add player account persistence or session persistence.
- Do not change the Protobuf envelope or WebSocket handshake.
- Do not add runtime player handlers or WebSocket routes.
- Do not add new game-domain modules.

## Acceptance Criteria

- `M-006` is completed.
- `M-007` is active.
- `W-0037` records the selected direction.
- Exactly one follow-up tooling work item is made ready or completed by subsequent tooling work.
