# Impact

## Summary

This is a source-first tooling and architecture-memory change. It changes repository inspection output, architecture manifests, check rules, and continuation metadata only.

It opens `W-0256 Define Pitaya-aligned cluster-safe session routing boundary gate` as the next bounded follow-up.

## Runtime Impact

No runtime behavior.

The current runtime remains a statically composed single-process server with target-scope metadata, application-owned server-push intent, and single-process WebSocket delivery.

## Protocol Impact

No protocol shape.

No WebSocket route, Protobuf message, envelope field, generated output, or protocol bridge changes are added.

No generated output.

## Persistence Impact

No persistence behavior.

No repository interface, PostgreSQL adapter, migration, table, group membership registry, room registry, stream subscription storage, or fanout persistence is added.

## Dependency Impact

No dependencies.

No group library, broadcast library, queue, stream, membership system, service discovery library, RPC library, load balancer, or observability backend is added.

## Security And Redaction Impact

No sensitive runtime state is exposed.

The inspection output keeps raw credentials, raw access tokens, lookup digests, verifier digests, verifier keys, DSNs, database payloads, local secret file contents, node credentials, and transport metadata out of output.

## Compatibility Impact

No direct Nakama/Pitaya API compatibility.

The map uses Pitaya as architecture vocabulary only. It does not add Pitaya package names, method names, group names, room names, route names, wire shapes, or behavior compatibility.
