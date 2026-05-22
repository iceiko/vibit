# Storage Objects Behavior Gate

Status: Accepted v0.1
Last updated: 2026-05-22
Scope: Gate for first general storage-object behavior beyond the inventory proof slice
Depends on: `docs/prototype-ready-local-development-path-package.md`, `docs/nakama-pitaya-product-parity-roadmap.md`, `docs/reference-game-server-alignment.md`
Canonical decision: `ADR-0109`

The paired Simplified Chinese translation is `docs/storage-objects-behavior-gate.zh-CN.md`. The English file is authoritative.

This document defines the first storage objects behavior gate. It is a gate artifact. It does not implement storage objects runtime behavior, add protocol routes, add Protobuf source or generated output, add migrations, add dependencies, add repository interfaces, add storage adapters, broaden operations/admin behavior, add hosted deployments, create release artifacts, run public announcements, run paid promotion, change authentication/session behavior, add a broad product module implementation, or add direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The storage objects behavior gate record is:

```yaml
storage_objects_behavior_gate: defined
completed_work_item: W-0201
decision: ADR-0109
check_rule: runtime.storage_objects_behavior_gate
source_package_decision: ADR-0108
source_package_standard: docs/prototype-ready-local-development-path-package.md
gate_standard: docs/storage-objects-behavior-gate.md
gate_standard_translation: docs/storage-objects-behavior-gate.zh-CN.md
target_stage: prototype_ready_foundation
reference_capability_family: storage_objects_and_durable_game_state
first_scope_posture: player_owned_small_json_objects
object_identity_tuple: owner_kind_owner_id_collection_key
ownership_posture_recorded: true
scope_key_posture_recorded: true
read_write_semantics_recorded: true
permission_posture_recorded: true
conflict_semantics_recorded: true
protocol_expectations_recorded: true
data_expectations_recorded: true
verification_expectations_recorded: true
stop_conditions_recorded: true
future_schema_gate_work_item: W-0202
future_schema_gate_direction: storage_objects_persistence_schema_gate
gate_only: true
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
dependency_added: false
repository_interface_changed: false
storage_adapter_changed: false
broad_operations_admin_behavior_added: false
authentication_session_behavior_changed: false
product_module_expansion_added: false
hosted_deployment_added: false
additional_release_artifacts_authorized: false
public_announcements_beyond_github_release_authorized: false
paid_promotion_authorized: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Product Intent

The inventory proof slice is useful, but it is module-specific. Prototype authors also need a general durable state surface for small records such as loadouts, preferences, tutorial state, quest flags, profile fragments, saved selections, or prototype-specific state that does not deserve a first-class module yet.

The first storage objects behavior should make vibit more useful for prototypes while preserving the agent-native model:

- ownership is explicit;
- object identity is stable;
- permissions are route-scoped and fail closed;
- writes are server-authoritative and version-aware;
- data shape is bounded and redaction-aware;
- persistence remains behind repository and adapter gates;
- protocol shape follows behavior rather than inventing wire messages first.

This gate adapts Nakama's broad durable storage-object capability as a product reference and Pitaya's handler/persistence separation as an architecture reference. It does not copy Nakama public routes, data models, permission values, version strings, or client API names, and it does not copy Pitaya routing or cluster internals.

## 3. Storage Object Boundary

In this gate, a storage object means a small durable game-state record owned by the server framework and addressed by a stable logical identity.

The first posture is:

```text
owner_kind: player
owner_id: authenticated player id
collection: server-defined logical namespace
key: server-defined or validated object key
value: bounded JSON object payload candidate
version: server-issued optimistic concurrency token
```

This is not large binary object storage, asset storage, replay storage, CDN-backed content, S3-compatible blob storage, file uploads, or arbitrary document database adoption. The existing `object_storage` planning entries in `.arch/runtime.yaml` remain separate and deferred.

## 4. Ownership And Scope

The first behavior gate selects player-owned objects only.

Allowed first scope:

- `player` owner kind;
- one authenticated player owns each first-posture object;
- the owner id must come from validated request identity, not from client-supplied metadata alone;
- collection and key identify the object within that player scope.

Deferred scopes:

- global objects;
- group, guild, party, room, match, or server-shard objects;
- cross-player shared writes;
- public catalog objects;
- admin-managed objects;
- object ACL lists;
- object search across arbitrary owners.

Those scopes require later gates because they change permissions, indexing, operations, and abuse/failure behavior.

## 5. Key And Value Posture

The first object identity is the tuple:

```text
owner_kind + owner_id + collection + key
```

The future implementation should treat `collection` and `key` as protocol-visible identifiers with explicit validation. The recommended first posture is ASCII-safe, length-bounded, case-sensitive strings with no path semantics.

The first value posture is a bounded JSON object payload candidate:

- the payload must be a JSON object, not arbitrary text or a binary blob;
- maximum size must be ratified before implementation;
- nested depth and field count should be bounded before production hardening;
- value contents are not log-safe by default;
- object ids, collection names, and keys may become log-safe only after a redaction decision.

The first value posture may be revised in the schema gate if the repository chooses a stricter format, but it must remain small-object game state, not general file or blob storage.

## 6. Read Semantics

The future first behavior should support route-scoped reads for the authenticated player's own objects.

Candidate read operations:

- get one object by collection and key;
- list objects in one collection for the authenticated player, with bounded pagination.

Read behavior should:

- require validated player identity;
- return object metadata and value only through protected routes;
- use stable public not-found behavior;
- avoid leaking whether another player's object exists;
- never trust client-supplied owner id over request identity in the first posture.

Cross-player reads, public reads, admin reads, and indexed search are deferred.

## 7. Write Semantics

The future first behavior should support server-authoritative writes for the authenticated player's own objects.

Candidate write operations:

- put or replace an object by collection and key;
- delete an object by collection and key.

Write behavior should:

- require validated player identity;
- validate collection, key, value shape, value size, and permission before persistence;
- issue a new version after each successful mutation;
- record create, update, and delete behavior in domain terms before protocol messages are added;
- avoid hidden partial writes;
- keep mutation logic inside the future storage objects module or application boundary, not transport or persistence adapters.

Multi-object transactions, batch writes, partial JSON patch, merge semantics, TTL, public ACL changes, and server-side script hooks are deferred.

## 8. Permission Posture

The first permission posture is fail-closed and route-scoped.

The first implementation gate should define permissions before code. Candidate permission families:

- read own storage object;
- list own storage objects in a collection;
- write own storage object;
- delete own storage object.

The first posture must not treat metadata-only `player_id` as proof. It must compose with the existing authenticated request path and the current route-protection model unless a later ADR explicitly changes authentication/session behavior.

Client-controlled read/write permission bits, public object permissions, object ACLs, and admin bypass are deferred.

## 9. Conflict Semantics

The first write posture should use optimistic concurrency.

Candidate version behavior:

- successful create or update returns a server-issued version;
- update and delete may accept an expected version;
- missing expected version behavior must be selected before implementation;
- stale expected version returns a stable conflict class;
- malformed version returns a validation class;
- not-found, permission failure, and owner mismatch must avoid cross-player existence leaks.

The schema gate should decide whether the first stored version is a monotonic integer, database revision, opaque server token, hash-derived token, or another explicit representation. The behavior gate only requires the version to be server-issued and not client-authoritative.

## 10. Protocol Expectations

Future protocol work should follow semantic behavior and contract boundaries.

Candidate route families are intentionally not implemented by this gate:

- `runtime.storage.GetObject`
- `runtime.storage.ListObjects`
- `runtime.storage.PutObject`
- `runtime.storage.DeleteObject`

The route names above are planning candidates only. A future protocol gate must define exact module/name routing, Protobuf source files, generated output, request/response shapes, error mapping, and compatibility posture before implementation.

The first protocol surface should preserve:

- WebSocket-framed Protobuf envelope use;
- protected-route authentication;
- application-owned dispatch and policy;
- transport credential-neutrality;
- generated-output traceability.

## 11. Data Expectations

The next bounded direction is a storage objects persistence schema gate.

That gate should decide:

- table name candidate, likely `storage_objects`;
- owner kind and owner id representation;
- collection/key constraints;
- value representation;
- version representation;
- create/update/delete timestamps;
- uniqueness and indexes;
- redaction posture for logs and diagnostics;
- migration source candidate;
- future repository and PostgreSQL adapter boundaries.

This behavior gate does not add SQL migration source, repository interfaces, PostgreSQL adapters, runtime composition, startup wiring, or migration apply behavior.

## 12. Verification Expectations

Future implementation gates should require focused tests for:

- create then read;
- replace with version change;
- stale version conflict;
- delete then not found;
- own-object permission success;
- cross-owner access rejection without existence leak;
- collection/key validation;
- value size and shape validation;
- redacted errors and logs;
- PostgreSQL persistence behavior after the schema and adapter are authorized.

Default repository checks must not require live PostgreSQL unless the check is explicitly opt-in through a disposable database environment.

## 13. Stop Conditions

Stop and ask for maintainer authorization before doing any of the following:

- implementing storage objects runtime behavior;
- adding protocol routes;
- adding Protobuf source files or generated output;
- adding SQL migrations;
- adding repository interfaces;
- adding PostgreSQL or other storage adapters;
- adding dependencies;
- changing authentication/session semantics;
- changing route protection semantics;
- adding cross-player, global, group, party, room, match, public, admin, or ACL object scopes;
- adding large object/blob storage or S3-compatible object storage;
- adding server-side custom logic hooks;
- adding broad operations/admin behavior;
- adding hosted deployments or demos;
- creating release binaries, packages, containers, checksums, signing/provenance artifacts, install scripts, registry publications, SDK packages, or additional release artifacts;
- executing public announcements beyond the GitHub release record;
- running paid promotion;
- adding direct Nakama/Pitaya API compatibility.

## 14. Next Work

The next bounded direction is:

```text
W-0202 Define storage objects persistence schema gate
```

That work should decide the first persistence schema posture for storage objects before any migration source, repository interface, adapter, protocol, or runtime implementation is added.
