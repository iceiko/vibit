# Plan

1. Add the English session repository boundary standard.
2. Add the Simplified Chinese translation.
3. Add `ADR-0061`.
4. Add the conversation log.
5. Update runtime, conventions, contracts, reference, work-item, module, and AGENTS manifests.
6. Add `runtime.session_repository_boundary` to rules and `tools/vibit`.
7. Run focused and full repository verification.

## Rollback

Reversal means removing the boundary standard, ADR, conversation log, check rule, manifest references, and `M-062/W-0134`, then reopening `M-061/W-0133` before any future repository implementation depends on this boundary.
