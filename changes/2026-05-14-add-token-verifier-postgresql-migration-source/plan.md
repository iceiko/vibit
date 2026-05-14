# Plan

1. Add `runtime/migrations/postgres/000004_create_authentication_access_tokens.sql`.
2. Update M-014 work state so `W-0078` is completed and `W-0079` becomes `next_ready`.
3. Update runtime, convention, contract, reference, persistence, selected login/token, schema gate, and agent guide artifacts to record token verifier migration source state.
4. Extend local checks to accept the token verifier migration as the authorized schema source while preserving runtime authentication deferral.
5. Run focused and full repository verification.
6. Commit and push the completed work item.
