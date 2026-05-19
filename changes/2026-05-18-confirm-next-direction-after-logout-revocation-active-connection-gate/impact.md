# Impact

This change closes the confirmation gate after the logout/revocation active-connection gate and opens a gate-only logout access-token behavior slice.

It does not add runtime code, Protobuf files, generated output, migrations, dependencies, WebSocket behavior, token revocation execution, runtime session revocation, active connection invalidation, or direct Nakama/Pitaya API compatibility.

Expected manifest impact:

- Mark `M-083/W-0155` completed.
- Add `M-084/W-0156` for `define_logout_access_token_behavior_gate`.
- Add `M-085/W-0157` as the next blocked confirmation gate.
- Update runtime, conventions, contracts, reference, module, AGENTS, rules, and checks to record the selected direction.
