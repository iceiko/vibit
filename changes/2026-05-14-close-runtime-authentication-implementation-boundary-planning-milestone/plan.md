# Plan

1. Review `M-016` against the runtime authentication implementation boundary standard, manifests, guides, checks, and work queue.
2. Add ADR-0037 to record the closeout and next gate.
3. Add a conversation log for the closeout.
4. Mark `M-016` and `W-0087` completed.
5. Open `M-017 Authentication Generated Contract Shape Gate`.
6. Mark `W-0088 Decide authentication generated contract shape timing` as next ready.
7. Update runtime, conventions, contracts, reference, and agent guides.
8. Run verification.

## Rollback

Revert the work-item and manifest updates to restore `M-016` as active and `W-0087` as next ready. No runtime or database rollback is needed because this change does not add code or migrations.
