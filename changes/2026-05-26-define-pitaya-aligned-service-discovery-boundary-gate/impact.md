# Impact

## Summary

This is a source-first architecture gate. It changes repository documentation, architecture manifests, check rules, and continuation metadata only.

It opens `W-0253 Implement Pitaya-aligned service discovery source-first map` as the next bounded follow-up.

## Runtime Impact

No runtime behavior.

The current runtime remains a statically composed single-process server with in-process dispatch and module handlers.

## Protocol Impact

No protocol shape.

No WebSocket route, Protobuf message, envelope field, generated output, or protocol bridge changes are added.

No generated output.

## Persistence Impact

No persistence behavior.

No repository interface, PostgreSQL adapter, migration, table, registry storage, node registry, or service registry persistence is added.

## Dependency Impact

No dependencies.

No service discovery library, registry client, RPC library, load balancer, membership system, or observability backend is added.

## Security And Redaction Impact

No sensitive runtime state is exposed.

The gate keeps raw credentials, raw access tokens, lookup digests, verifier digests, verifier keys, DSNs, database payloads, local secret file contents, node credentials, and transport metadata out of artifacts and inspection outputs.

## Compatibility Impact

No direct Nakama/Pitaya API compatibility.

The gate uses Pitaya as architecture vocabulary only. It does not add Pitaya package names, method names, registry names, route names, wire shapes, or behavior compatibility.
