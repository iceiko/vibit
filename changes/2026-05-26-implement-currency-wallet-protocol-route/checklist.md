# Checklist

- [x] Add currency wallet route keys and route kind tests.
- [x] Add `proto/vibit/currency/v1/currency.proto`.
- [x] Generate `runtime/internal/generated/proto/vibit/currency/v1/currency.pb.go` through Buf.
- [x] Add currency protocol bridge request/response mapping.
- [x] Add currency bootstrap route handlers.
- [x] Require a server-side grant policy before `GrantCurrency` route execution.
- [x] Keep validated authenticated request identity as owner proof.
- [x] Omit owner/player/session/wallet lookup/actor proof fields from currency Protobuf payloads.
- [x] Preserve metadata-only identity refusal through the existing protected route path.
- [x] Add route, bridge, handler, and local authenticated gameplay tests.
- [x] Allow zero-field proto3 authenticated inner payloads without weakening non-empty payload validation.
- [x] Add `runtime.currency_wallet_protocol_route_implementation`.
- [x] Update work items, manifests, guides, README/docs, ADR, change spec, and conversation log.
- [x] Record that no next-ready work item currently exists after W-0302.
- [x] Preserve reward, inventory, purchase, catalog, event, payment, migration, dependency, startup, hosted, SDK, distributed runtime, and direct compatibility deferrals.
