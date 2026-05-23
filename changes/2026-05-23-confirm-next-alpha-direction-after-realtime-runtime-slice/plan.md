# Plan

1. Read the work queue, runtime manifests, first server push/realtime gate, realtime runtime slice ADR, product maturity milestones, and Nakama/Pitaya roadmap.
2. Select the next bounded alpha direction after the realtime runtime slice.
3. Record the decision in an ADR and conversation memory.
4. Update the work queue, runtime/reference manifests, README, alpha docs, AGENTS guides, and repository checks.
5. Verify repository checks and record results.

## Non-Implementation Rule

This is a direction-selection slice only. Do not implement WebSocket outbound delivery, socket writes, realtime protocol bridges, Protobuf messages, generated output, startup wiring, persistence, stream subscriptions, chat, groups, broadcast fanout, authentication/session changes, migrations, dependencies, delivery guarantees, distributed runtime, match runtime, or direct Nakama/Pitaya API compatibility in this change.
