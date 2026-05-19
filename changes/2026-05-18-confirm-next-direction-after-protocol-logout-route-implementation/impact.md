# Impact

This change records a work-queue decision only.

It completes `M-099/W-0171` and opens `M-100/W-0172` as the next gate-only lifecycle slice:

```text
define_transport_close_handoff_gate
```

## Affected Scope

- Runtime lifecycle planning.
- Work queue state.
- Agent handoff context.
- Product parity sequencing.

## No Runtime Behavior Added

This confirmation step does not add:

- Go WebSocket close handoff code.
- Socket close codes.
- Player-facing close reason text.
- Protocol close messages.
- Protobuf messages or generated output.
- Logout-triggered socket close.
- Runtime session revocation.
- Reconnect/epoch behavior.
- Protocol session carriers.
- Presence, chat, social, matchmaking, match runtime, SDK, cluster, or distributed runtime behavior.
- Dependencies.
- Direct Nakama/Pitaya API compatibility.

## Reference Alignment

The selected direction follows the roadmap's lifecycle closure order:

```text
protocol logout route -> transport close handoff -> reconnect/epoch -> protocol session carrier -> presence
```

This absorbs Nakama's lifecycle product pressure and Pitaya's connection-management layering without copying either public API surface.
