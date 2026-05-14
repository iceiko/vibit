# Plan

1. Record the maintainer's selected direction in a change spec and conversation log.
2. Mark `M-004` and `W-0030` completed.
3. Add `M-005 Player Account And Session Contracts` as the active milestone.
4. Add exactly one `next_ready` item for the first contract standard step.
5. Update runtime and reference manifests to point agents at the selected direction.
6. Run work, change, and repository checks.
7. Update the verification record.

## Boundaries

This plan must not:

- Choose an authentication scheme.
- Choose a token format.
- Choose credential storage or password behavior.
- Add player account migrations.
- Add session persistence.
- Change the Protobuf envelope.
- Change the WebSocket handshake.
- Copy Nakama or Pitaya public API shape.
