# Impact Analysis

## Affected Modules

- `runtime`: adds the first frame-to-application composition adapter under the Protobuf protocol adapter package.
- `inventory`: used as the first domain route family in composition tests; no inventory contract or domain behavior changes.

## Module Ownership Impact

No module ownership changes.

The new adapter lives under `runtime/internal/platform/protocol/protobuf/` because it owns Protobuf envelope decoding, generated Protobuf imports, protocol/domain payload bridging, and application-result encoding. It does not import the WebSocket transport package, keeping `runtime/internal/platform/transport/ws/` narrow and opaque.

## Public Contract Impact

No command, query, event, error, permission, or Protobuf envelope shape changes.

## Data And Migration Impact

No data or migration impact. Persistence remains deferred.

## Test Impact

Adds protocol adapter tests for:

- successful GrantItem frame handling,
- application error frame handling,
- malformed frame rejection.

## Documentation Impact

Updates runtime protocol adapter documentation, runtime manifest state, runtime agent guide state, and the work-item queue.

## Compatibility Risks

The adapter returns protocol decode and internal non-application errors as Go errors rather than inventing a public protocol error shape. This preserves the current protocol decision and avoids leaking internal failures.
