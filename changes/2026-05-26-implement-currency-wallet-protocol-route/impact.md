# Impact

## Runtime

This change adds a protected currency wallet route family to the existing application dispatch and Protobuf adapter layers.

The application service remains owned by `runtime/internal/app/currency`. The route implementation does not move wallet behavior into WebSocket transport, generated Protobuf code, PostgreSQL adapters, or module value types.

## Protocol

This change adds `vibit.currency.v1` Protobuf request and response messages and generated Go output. The envelope shape is unchanged.

`EnsurePlayerWalletRequest` and `GetOwnWalletRequest` are zero-field proto3 messages. The authenticated frame handler now permits empty encoded inner payload bytes only when the registered Protobuf descriptor for the inner payload type has zero fields. Empty payloads for non-empty messages remain malformed.

## Identity And Authorization

Validated authenticated request identity remains the owner proof source. Client payloads do not carry owner ids, player ids, session ids, wallet lookup ids, actor ids, access tokens, credential material, lookup digests, verifier digests, or payment provider payloads.

`GrantCurrency` requires an explicit server-side `CurrencyGrantPolicy`. The route handler refuses grant execution when the policy is absent, denied, or lacks a system actor id.

## Data

No migrations are added or changed. The implementation uses existing currency wallet tables through the existing application service and repository boundaries in tests.

## Compatibility

The route names are vibit-native and do not add direct Nakama or Pitaya public API compatibility. Broader economy features remain deferred.
