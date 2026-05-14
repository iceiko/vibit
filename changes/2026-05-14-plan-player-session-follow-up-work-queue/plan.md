# Plan

1. Inspect the current M-003 milestone and next-ready work item.
2. Add the remaining conservative M-003 work items after `W-0025`.
3. Mark one follow-up work item as `next_ready`.
4. Keep authentication, token, credential, player account persistence, Protobuf envelope, and WebSocket handshake decisions behind ask-first boundaries.
5. Run work inspection and repository checks.
6. Record verification.

## Resulting Queue Shape

- `W-0025`: Plan player/session follow-up work queue.
- `W-0026`: Add session validator hook boundary.
- `W-0027`: Add inventory request identity permission handoff boundary.
- `W-0028`: Add metadata-only identity repository checks.
- `W-0029`: Close M-003 player identity and session boundary.

This order gives future agents a small implementation hook before moving inventory away from bootstrap-only permission assumptions, then adds checks to prevent regression before the milestone is closed.
