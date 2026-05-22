# Conversation: Storage Objects Persistence Schema Gate

Date: 2026-05-22
Participants: Maintainer, Agent
Related work item: `W-0202`
Related decision: `ADR-0110`

## Context

The maintainer asked to continue advancing vibit after the storage objects behavior gate was completed. The next-ready work item was `W-0202 Define storage objects persistence schema gate`.

The behavior gate selected player-owned small JSON storage objects as the first durable game-state posture and deferred schema, migration, repository, adapter, protocol, and runtime behavior work.

## Maintainer Narrative

The maintainer wants vibit to keep moving toward a real prototype-ready and eventually production-useful server framework. A storage object schema gate is the next concrete step because it turns the behavior gate into an inspectable persistence plan without yet adding SQL or runtime behavior.

## Agent Response Summary

The agent defined the storage objects persistence schema gate as a planning-only, schema-gate-only slice. The gate selects PostgreSQL, a future `storage_objects` table, the future migration source candidate `runtime/migrations/postgres/000006_create_storage_objects.sql`, JSONB values, BIGINT versions, player-owner identity, collection/key constraints, timestamp and soft-delete posture, uniqueness/index posture, redaction rules, and future repository/adapter boundaries.

## Decisions

- The storage objects persistence schema gate is recorded in `docs/storage-objects-persistence-schema-gate.md`.
- The paired Simplified Chinese translation is `docs/storage-objects-persistence-schema-gate.zh-CN.md`.
- The decision record is `ADR-0110`.
- The repository check rule is `runtime.storage_objects_persistence_schema_gate`.
- The future logical table candidate is `storage_objects`.
- The future migration source candidate is `runtime/migrations/postgres/000006_create_storage_objects.sql`.
- The next bounded direction is `W-0203 Add storage objects migration source`.

## Artifacts

- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`
- `.arch/reference.yaml`
- `docs/storage-objects-persistence-schema-gate.md`
- `docs/storage-objects-persistence-schema-gate.zh-CN.md`
- `decisions/ADR-0110-storage-objects-persistence-schema-gate.md`
- `changes/2026-05-22-define-storage-objects-persistence-schema-gate/`
- `rules/check-rules.json`
- `tools/vibit`

## Open Questions

- The actual SQL migration source is deferred to `W-0203`.
- The storage repository interface remains deferred.
- The PostgreSQL adapter remains deferred.
- Protocol routes, Protobuf source, generated output, and runtime behavior remain deferred.

## Follow-Up

The next work item may add `runtime/migrations/postgres/000006_create_storage_objects.sql` as a migration-source-only slice with matching static checks.

## Redaction Notes

No private local environment file was read or printed. No secrets, raw credentials, raw tokens, verifier keys, DSNs, digests, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, object values, or concrete transport metadata were added to the record.
