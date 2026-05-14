# Impact

## Affected Modules

- `runtime/internal/app`
- `runtime/cmd/vibit-server`

The `inventory` and `player` modules are not changed by this step.

## Module Ownership Impact

Application dispatch becomes the owner of a concrete session validation hook boundary. This matches `ADR-0021`: validation belongs after protocol decoding and before module handlers receive route requests.

Transport still owns only WebSocket connection mechanics. Protobuf still owns only envelope decoding and metadata conversion. Inventory still owns inventory behavior and permission policy, not session validation or player accounts.

## Public Contract Impact

No public commands, queries, events, errors, permissions, Protobuf messages, or WebSocket routes are added or changed.

## Runtime Impact

The runtime request path can now wrap a dispatcher with `SessionValidatingDispatcher`. The default validator produces or preserves metadata-only identity and does not authenticate clients.

An injected validator may return `Valid: true` with a replacement identity, or `Valid: false` to stop dispatch with an application error. This is hook behavior only; no concrete validator is implemented.

## Protocol Impact

No Protobuf envelope fields change. No WebSocket handshake behavior changes.

## Data And Migration Impact

No database schema, migration, session store, credential store, or token store is added.

## Test Impact

Focused Go tests cover:

- Metadata-only validator normalization.
- Session-validating dispatcher pass-through behavior.
- Injected validated identity replacement.
- Invalid validation preventing handler execution.
- Missing dispatcher and missing validator errors.

## Documentation Impact

Runtime and player/session boundary docs are updated in English and Simplified Chinese to record the hook without claiming real authentication exists.

## Compatibility Risks

The hook is internal Go runtime behavior. Wire compatibility is unchanged. The main risk is accidentally interpreting metadata-only identity as authenticated proof; tests and documentation keep validation flags false for the default path.
