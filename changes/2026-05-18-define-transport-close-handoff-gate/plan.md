# Plan

1. Add the English transport close handoff gate standard.
2. Add the Simplified Chinese translation.
3. Add `ADR-0080`.
4. Record the gate in architecture manifests and AGENTS guides.
5. Add repository check metadata for `runtime.transport_close_handoff_gate`.
6. Complete `M-100/W-0172`.
7. Create a next confirmation point for the implementation direction.
8. Verify work queue, change spec, runtime checks, and next inspection.

## Non-Implementation Boundary

This plan intentionally does not edit WebSocket transport code, close policy code, protocol bridge behavior, authentication service behavior, Protobuf sources, generated output, migrations, or dependencies.
