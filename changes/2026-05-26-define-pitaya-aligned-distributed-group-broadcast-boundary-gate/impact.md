# Impact

## Summary

This is a source-first architecture gate. It changes repository documentation, architecture manifests, check rules, and continuation metadata only.

It opens `W-0255 Implement Pitaya-aligned distributed group and broadcast source-first map` as the next bounded follow-up.

## Runtime Impact

No runtime behavior.

The current runtime remains single-process. Target-scope metadata and application-owned server-push intent are mapped as planning vocabulary only.

## Protocol Impact

No protocol shape.

No WebSocket route, Protobuf message, envelope field, generated output, or protocol bridge changes are added.

No generated output.

## Persistence Impact

No persistence behavior.

No repository interface, PostgreSQL adapter, migration, table, group membership registry, stream subscription table, room state, fanout queue, or broadcast persistence is added.

## Dependency Impact

No dependencies.

No group, broadcast, stream, queue, service discovery, RPC, load balancer, membership, or observability backend library is added.

## Security And Redaction Impact

No sensitive runtime state is exposed.

The gate keeps raw credentials, raw access tokens, lookup digests, verifier digests, verifier keys, DSNs, database payloads, local secret file contents, node credentials, transport metadata, group membership payloads, and broadcast payloads out of artifacts and inspection outputs.

## Compatibility Impact

No direct Nakama/Pitaya API compatibility.

The gate uses Pitaya as architecture vocabulary only. It does not add Pitaya or Nakama package names, method names, group names, room names, stream names, route names, wire shapes, or behavior compatibility.
