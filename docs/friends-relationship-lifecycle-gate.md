# Friends Relationship Lifecycle Gate

Status: Accepted v0.1
Last updated: 2026-05-24
Scope: Gate for future player friendship relationship lifecycle semantics before persistence, protocol, runtime behavior, or broader social features
Depends on: `docs/agent-native-feature-request-test-workflow.md`, `docs/agent-native-feature-request-scaffolding.md`, `docs/reference-game-server-alignment.md`, `docs/nakama-pitaya-product-parity-roadmap.md`
Canonical decision: `ADR-0139`

The paired Simplified Chinese translation is `docs/friends-relationship-lifecycle-gate.zh-CN.md`. The English file is authoritative.

This document defines the friends relationship lifecycle semantic gate. It is a gate artifact. It does not add runtime friendship behavior, protocol routes, Protobuf source, generated output, migrations, repository interfaces, PostgreSQL adapters, dependencies, startup wiring, authentication/session behavior changes, delivery guarantees, stream subscriptions, chat rooms, groups, parties, broadcast fanout, matchmaking, match runtime, operations/admin behavior, SDK publication, generated client libraries, hosted deployments, release artifacts, Pitaya-style distributed architecture, or direct Nakama/Pitaya API compatibility.

## 1. Core Rule

The friends relationship lifecycle gate record is:

```yaml
friends_relationship_lifecycle_gate: defined
completed_work_item: W-0231
decision: ADR-0139
check_rule: runtime.friends_relationship_lifecycle_gate
source_intake_decision: ADR-0138
source_scaffolding_decision: ADR-0137
source_workflow_decision: ADR-0128
gate_standard: docs/friends-relationship-lifecycle-gate.md
gate_standard_translation: docs/friends-relationship-lifecycle-gate.zh-CN.md
selected_nakama_capability_family: friends_groups_and_parties
primary_product_reference: Nakama
pitaya_reference_status: deferred_future_architecture_reference
ai_native_development_testing_goal: user_requirement_to_spec_tests_implementation_verification
semantic_gate_only: true
future_persistence_schema_gate_work_item: W-0232
future_persistence_schema_gate_direction: define_friends_relationship_persistence_schema_gate
runtime_behavior_added: false
protocol_route_added: false
protobuf_source_added: false
generated_output_added: false
migration_added: false
repository_interface_changed: false
postgresql_adapter_changed: false
dependency_added: false
direct_nakama_pitaya_api_compatibility_added: false
```

## 2. Product Intent

Friendship is a core social graph primitive in Nakama-class game backends. It can support later groups, parties, chat targeting, invites, matchmaking filters, player discovery, and match social context. vibit adopts the product capability, not the public API shape.

The vibit posture is:

- contract-first lifecycle semantics before implementation;
- server-authoritative state transitions;
- actor-relative public relationship status;
- validated player identity for every future command and query;
- explicit test expectations before code;
- private social graph data treated as not log-safe by default;
- direct external API compatibility rejected unless a later ADR authorizes it.

Pitaya remains deferred as a future distributed architecture reference. This gate must not pull in distributed routing, frontend/backend roles, RPC, cluster groups, service discovery, or server-to-server messaging.

## 3. Future Semantic Scope

The future friends relationship lifecycle must cover:

```yaml
semantic_scope:
  - request
  - accept
  - reject
  - remove
  - block
  - unblock
  - list
  - read_relationship_status
```

The lifecycle is player-to-player and server-authoritative. The future domain owner is a social/friends capability boundary, not WebSocket transport, protocol adapters, authentication, storage objects, inventory, realtime delivery, chat, matchmaking, or match runtime.

## 4. Future Contract Vocabulary

The future command vocabulary is:

```yaml
commands:
  - SendFriendRequest
  - AcceptFriendRequest
  - RejectFriendRequest
  - RemoveFriend
  - BlockPlayer
  - UnblockPlayer
```

The future query vocabulary is:

```yaml
queries:
  - ListFriendRelationships
  - GetFriendRelationshipStatus
```

The future event vocabulary is:

```yaml
events:
  - FriendRequestCreated
  - FriendRequestAccepted
  - FriendRequestRejected
  - FriendRemoved
  - PlayerBlocked
  - PlayerUnblocked
```

The future error vocabulary is:

```yaml
errors:
  - FRIENDSHIP_INVALID_TARGET
  - FRIENDSHIP_SELF_TARGET_FORBIDDEN
  - FRIENDSHIP_DUPLICATE_REQUEST
  - FRIENDSHIP_BLOCKED_RELATIONSHIP
  - FRIENDSHIP_INVALID_TRANSITION
  - FRIENDSHIP_RELATIONSHIP_NOT_FOUND
  - FRIENDSHIP_TARGET_NOT_FOUND
  - FRIENDSHIP_METADATA_IDENTITY_NOT_AUTHENTICATED
```

This vocabulary is semantic planning. It does not create contract source files, generated shapes, protocol payloads, routes, repositories, or runtime handlers.

## 5. Identity And Permissions

Every future command and query requires:

```yaml
permission: validated_player_identity
metadata_only_player_id_allowed_as_proof: false
metadata_only_session_id_allowed_as_proof: false
actor_identity_source: validated_request_identity
```

Rules:

- The actor is derived from validated request identity, not from a client-supplied actor id.
- `player_id` and `session_id` metadata are not authentication proof.
- Self-targeting is forbidden.
- Target existence and hidden relationship details must not leak through public failures when authorization or privacy requires collapse.
- Relationship records, target ids, actor ids, private statuses, and conflict details are not log-safe by default.

## 6. Relationship State Model

The future public relationship status is actor-relative:

```yaml
actor_relative_public_states:
  - none
  - outgoing_request_pending
  - incoming_request_pending
  - friends
  - blocked_by_actor
  - blocked_actor
  - mutual_blocked
  - removed
  - rejected
```

The future persisted canonical state may be pair-oriented. If it is pair-oriented, public query output must still be computed relative to the requesting actor. The storage shape is intentionally deferred to `W-0232`.

The first lifecycle invariants are:

- Self-targeting is forbidden.
- Duplicate request behavior must be either forbidden or explicitly idempotent before implementation.
- Accept applies only to an incoming pending request.
- Reject applies only to an incoming pending request.
- Remove applies only to an existing friendship and must not remove a block.
- Block removes or overrides pending and friend relationships in favor of the actor's block.
- Unblock removes the actor's block and must not automatically restore a prior friendship.
- Mutual block status is represented only when both sides have an active block.
- List and status queries are actor-relative and must not expose another player's private graph.

## 7. Future Persistence Gate

The next bounded work item is:

```text
W-0232 Define friends relationship persistence schema gate
```

That follow-up should define table candidates, pair identity posture, block representation, indexes, uniqueness, event/audit posture, redaction, and repository/adapter boundaries before migration source exists.

This gate intentionally does not decide:

- exact table names;
- canonical pair key encoding;
- whether rejected or removed states are tombstones, audit-only facts, or current rows;
- duplicate request idempotency;
- concurrency conflict resolution;
- hard delete or retention policy;
- repository interface shape;
- PostgreSQL adapter SQL;
- protocol routes or payloads.

## 8. Future Test Expectations

Future behavior tests must be planned before implementation.

Positive tests:

- send friend request;
- accept incoming request;
- reject incoming request;
- remove existing friend;
- block target player;
- unblock previously blocked player;
- list actor-relative relationships;
- read actor-relative relationship status.

Negative tests:

- self-targeting;
- duplicate request or explicitly chosen idempotency behavior;
- invalid transition;
- blocked relationship interaction;
- missing or unknown target;
- missing relationship;
- metadata-only identity.

Permission and authentication tests:

- every command/query requires validated player identity;
- client-supplied actor id is ignored or rejected;
- metadata-only `player_id` and `session_id` are rejected as proof.

Persistence and transaction tests:

- schema and repository tests must be defined after `W-0232` and before migration/adapter/runtime implementation;
- command transitions must be transactional;
- emitted events and state changes must be consistent within the future unit-of-work boundary.

Failure and redaction tests:

- public errors do not leak private relationship graph details where privacy requires collapse;
- logs do not expose raw credentials, tokens, verifier keys, digests, transport metadata, or private social graph internals.

Concurrency tests:

- simultaneous request, accept, reject, remove, block, and unblock conflicts must have explicit expected outcomes before runtime implementation.

Integration and end-to-end tests:

- deferred until protocol routes and runtime handlers are authorized.

## 9. Non-Authorization

This gate does not authorize:

- runtime friendship behavior;
- protocol routes;
- Protobuf source;
- generated output;
- migrations;
- repository interfaces;
- PostgreSQL adapters;
- dependencies;
- startup wiring;
- authentication/session behavior changes;
- delivery guarantees;
- stream subscriptions;
- chat rooms;
- groups;
- parties;
- broadcast fanout;
- matchmaking;
- match runtime;
- operations/admin behavior;
- SDK publication;
- generated client libraries;
- hosted deployments;
- release artifacts;
- Pitaya-style distributed architecture;
- direct Nakama/Pitaya API compatibility.

Any future work in those areas requires a separate bounded work item and verification record.

## 10. Verification

This gate is verified by:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.friends_relationship_lifecycle_gate
node tools/vibit check change define-friends-relationship-lifecycle-gate --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```
