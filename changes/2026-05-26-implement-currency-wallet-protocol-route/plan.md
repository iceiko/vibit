# Plan

1. Add a focused route-key test for the currency wallet route family.
2. Add `vibit.currency.v1` Protobuf source for wallet, balance, transaction, request, and response messages.
3. Run Buf generation so generated Go output is produced from the `.proto` source.
4. Add currency route key helpers under `runtime/internal/app/currency`.
5. Add currency bootstrap route handlers that inject validated request identity and require server-side grant policy for grants.
6. Add protocol bridge mapping between `vibit.currency.v1` payloads and application currency requests/results.
7. Register generated currency payloads in the Protobuf registry and dispatch bridge.
8. Add focused bridge and bootstrap tests for mapping, identity handoff, redaction, forbidden protocol fields, and grant policy refusal.
9. Add an authenticated gameplay E2E proof for ensure/get/grant/spend/list balances/list transactions and post-logout protected route failure.
10. Fix zero-field proto3 authenticated payload handling without allowing empty payloads for non-empty messages.
11. Add `ADR-0210`, change spec artifacts, conversation log, manifests, module guides, docs, and repository check coverage.
12. Run focused and full verification.
