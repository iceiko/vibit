# Plan

1. Add an English ratification document for the first login-method set.
2. Add the Simplified Chinese translation.
3. Add ADR-0025 for the durable login-method decision.
4. Add a conversation log preserving maintainer intent and agent summary.
5. Update architecture manifests to record the ratified set and deferred implementation status.
6. Mark W-0066 completed and promote W-0067 to next_ready.
7. Run repository verification.
8. Record verification results.

Rollback:

- Revert the W-0066 change if the ratified set needs to be superseded before implementation.
- If the method changes later, add a superseding ADR rather than silently editing the decision away.
