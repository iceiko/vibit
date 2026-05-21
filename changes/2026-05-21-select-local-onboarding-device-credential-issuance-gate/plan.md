# Plan

1. Mark `M-108/W-0180` completed with selected direction `define_local_onboarding_device_credential_issuance_gate`.
2. Create `M-109/W-0181` as the next ready gate-only work item.
3. Record `ADR-0088` for the next alpha direction selection.
4. Add conversation memory for the continuation step.
5. Update `.arch/work-items.yaml`, `.arch/runtime.yaml`, `.arch/conventions.yaml`, `.arch/contracts.yaml`, and `.arch/reference.yaml`.
6. Update alpha goal and agent-facing guidance so future sessions see the new queue state.
7. Update `tools/vibit` lifecycle phase expectations so historical runtime checks accept the queue advancing beyond `next_alpha_direction_selection`.
8. Verify the change spec, work queue, schemas, runtime checks, full repository checks, and diff whitespace.
