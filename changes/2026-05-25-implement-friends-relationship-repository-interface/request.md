# Request

Implement the storage-neutral friends relationship repository interface for the Nakama-first friends/groups/parties capability.

The implementation must stay inside the repository-interface slice:

- add module-owned Go repository vocabulary under `runtime/internal/modules/friends`;
- preserve canonical unordered pair identity, validated actor handoff posture, lifecycle state vocabulary, version/conflict semantics, and redacted errors from `ADR-0142`;
- add focused tests for normalization, copying, redaction, and forbidden dependency posture;
- add module manifest and module AGENTS guidance;
- keep PostgreSQL adapter behavior, SQL execution, runtime friendship behavior, protocol routes, Protobuf source, generated output, dependencies, migration changes, event/audit tables, SDKs, hosted deployment, distributed runtime, and direct Nakama/Pitaya API compatibility deferred.

## User-Facing Need

The project is aiming at a Nakama-class game/backend server framework shaped for AI-native development and testing. Friends relationships are a core social graph primitive, but future runtime behavior should not depend on ad hoc SQL-shaped or transport-shaped vocabulary.

This slice gives future agents a typed, storage-neutral repository contract before adapter and runtime behavior are authorized.

## Stop Conditions

- PostgreSQL adapter implementation or SQL execution is required.
- Runtime friend request/list/status behavior is required.
- Protocol routes, Protobuf source, or generated output are required.
- Authentication/session behavior, request identity validation, or startup wiring changes are required.
- Event/audit tables, chat, groups, parties, matchmaking, match runtime, SDKs, hosted deployment, distributed runtime, or direct compatibility are required.
- Verification fails.
