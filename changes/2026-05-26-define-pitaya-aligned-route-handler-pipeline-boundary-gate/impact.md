# Impact

## Runtime

No runtime behavior is added.

The current WebSocket Protobuf path remains:

```text
protocol envelope -> route request -> application dispatch -> command/query handler -> application result
```

The new vocabulary maps this current flow to deferred Pitaya-aligned route handler pipeline planning terms only.

## Protocol

No protocol shape is changed.

No protocol messages or routes, Protobuf sources, generated output, payload registry behavior, serializer behavior, or protocol bridge behavior are added.

No generated output.

## Persistence And Dependencies

No repository interfaces, PostgreSQL adapters, migrations, persistence behavior, or dependencies are added.

No dependencies.

## Product And Reference Alignment

Nakama remains the primary product capability reference. Pitaya remains a deferred architecture vocabulary reference for route handlers, handler pipelines, serializers, forwarding, frontend/backend roles, RPC, service discovery, groups, broadcast, and cluster routing.

## Follow-Up

`W-0260 Implement Pitaya-aligned route handler pipeline source-first map` becomes next-ready.
