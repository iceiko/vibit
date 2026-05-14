# Checklist

- [x] Confirm `W-0024` is the current next-ready work item.
- [x] Read the player identity and session boundary standard.
- [x] Add request identity and session-validation handoff types under `runtime/internal/app`.
- [x] Add metadata-only identity normalization helpers.
- [x] Add request identity to application route requests and results.
- [x] Preserve Protobuf envelope and WebSocket handshake shape.
- [x] Convert existing envelope session metadata into metadata-only request identity.
- [x] Preserve application identity in dispatcher, bootstrap, and inventory results.
- [x] Document the application identity handoff in runtime and player/session standards.
- [x] Record handoff status in architecture manifests.
- [x] Add focused Go tests.
- [x] Run verification.
- [x] Update `.arch/work-items.yaml`.
