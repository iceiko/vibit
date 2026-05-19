# Plan

1. Record direction confirmation from `M-055/W-0127` to `M-056/W-0128`.
2. Extend `proto/vibit/authentication/v1/authentication.proto` with BindConnection messages and status enum.
3. Run `buf generate` to regenerate Go Protobuf output.
4. Add application-owned connection binding types and binder in `runtime/internal/app`.
5. Add application binder tests for success, failure mapping, redaction, and metadata-only rejection.
6. Add Protobuf adapter bridge for `runtime.authentication.BindConnection`.
7. Route `BindConnection` through `FrameHandler` before normal route protection and dispatch.
8. Preserve the existing request-level `AuthenticatedRequest` protected-route path.
9. Add WebSocket connection epoch metadata handoff without parsing credential carriers.
10. Wire PostgreSQL startup composition with a connection binder when the authentication service is composed.
11. Update work queue, manifests, AGENTS guides, rules, tools, and conversation memory.
12. Run repository checks and Go tests.

## Rollback Or Migration Notes

No migration rollback is needed because no schema or durable state is added.
