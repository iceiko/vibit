# Plan

1. Add the English runtime authentication implementation boundary standard and Simplified Chinese translation.
2. Add `ADR-0036` for the boundary decision.
3. Add a conversation log for the planning decision.
4. Update runtime, conventions, contracts, reference, authentication module, and agent guides.
5. Mark `W-0086` completed and make the next conservative boundary work item ready.
6. Run focused and full verification.

## Rollback

Remove the boundary standard, translation, ADR, conversation log, manifest references, guide references, and work-item state changes. No code, migrations, or database rollback is needed because this change does not implement runtime behavior.
