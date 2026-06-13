# Request

Implement `W-0302 Implement currency wallet protocol route`.

The implementation must expose only the protected currency wallet route family authorized by `ADR-0209`:

- `currency.EnsurePlayerWallet`
- `currency.GetOwnWallet`
- `currency.ListOwnWalletBalances`
- `currency.GrantCurrency`
- `currency.SpendCurrency`
- `currency.ListOwnWalletTransactions`

The route family must use the existing WebSocket/Protobuf request path, generated Protobuf output, explicit protocol bridge mapping, application-owned handlers, validated authenticated request identity handoff, and a server-side grant policy.

## Scope

Allowed:

- Add `proto/vibit/currency/v1/currency.proto`.
- Generate `runtime/internal/generated/proto/vibit/currency/v1/currency.pb.go` through Buf.
- Add currency route keys under `runtime/internal/app/currency`.
- Add bootstrap handlers under `runtime/internal/app/bootstrap`.
- Add protocol bridge mapping under `runtime/internal/platform/protocol/protobuf`.
- Add focused route, bridge, handler, and authenticated local flow tests.
- Allow authenticated wrapper payloads with zero encoded bytes only when the inner Protobuf message descriptor has zero fields.

Not allowed:

- Reward integration.
- Inventory integration.
- Purchase behavior.
- Catalog, event, audit, or payment behavior.
- Reservation, settlement, refund, transfer, operations, or admin behavior.
- Migration changes.
- Dependency additions.
- Startup wiring.
- Authentication/session semantic changes.
- Hosted deployments, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.
