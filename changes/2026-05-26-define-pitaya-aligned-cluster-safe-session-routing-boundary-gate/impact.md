# Impact

## Summary

This is a source-first architecture gate. It changes repository documentation, architecture manifests, check rules, and continuation metadata only.

It opens `W-0257 Implement Pitaya-aligned cluster-safe session routing source-first map` as the next bounded follow-up.

## Runtime Impact

No runtime behavior.

The current runtime remains single-process. Server-observed connection id and epoch metadata, first-message connection binding, active connection registry vocabulary, request-level token identity, and runtime session validation are mapped as planning vocabulary only.

## Protocol Impact

No protocol shape.

No WebSocket route, Protobuf message, envelope field, generated output, proof carrier, session carrier, reconnect carrier, handoff message, or protocol bridge changes are added.

No generated output.

## Persistence Impact

No persistence behavior.

No repository interface, PostgreSQL adapter, migration, table, session location registry, connection owner node registry, routing cache, reconnect route table, service registry, or handoff state is added.

## Dependency Impact

No dependencies.

No session routing, service discovery, RPC, cluster runtime, load balancer, routing cache, membership, queue, or observability backend library is added.

## Security And Redaction Impact

No sensitive runtime state is exposed.

The gate keeps raw credentials, raw access tokens, lookup digests, verifier digests, verifier keys, DSNs, database payloads, local secret file contents, node credentials, transport metadata, connection payloads, session payloads, and routing metadata contents out of artifacts and inspection outputs.

## Compatibility Impact

No direct Nakama/Pitaya API compatibility.

The gate uses Pitaya as architecture vocabulary only. It does not add Pitaya or Nakama package names, method names, session names, route names, registry names, selector names, handoff names, wire shapes, or behavior compatibility.
