# Impact

This is a direction-confirmation record only.

It selects `implement_transport_close_handoff_single_process` as the next bounded lifecycle-closure step because `ADR-0080` already defined server-observed `connection_id + epoch` as the future concrete handoff target.

This change adds no runtime behavior, protocol behavior, migration, dependency, generated output, concrete WebSocket close, close code mapping, close reason text, logout-triggered socket close, runtime session revocation, reconnect behavior, protocol session carrier, operations/admin behavior, broad product module, or direct Nakama/Pitaya API compatibility.
