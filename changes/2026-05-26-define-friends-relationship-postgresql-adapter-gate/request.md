# Request

Define the friends relationship PostgreSQL adapter gate for the Nakama-first friends/groups/parties capability.

The gate must stay inside the adapter-boundary slice:

- define the future PostgreSQL adapter owner, source candidate, test candidate, and constructor posture;
- preserve the storage-neutral `runtime/internal/modules/friends.Repository` interface as the adapter contract;
- map the future SQL posture to the existing `friend_relationships` migration source without adding SQL execution behavior;
- define transaction/unit-of-work handoff, conflict mapping, redaction, and test expectations;
- open a bounded follow-up for implementation;
- keep actual PostgreSQL adapter code, runtime friendship behavior, protocol routes, Protobuf source, generated output, dependencies, migration changes, event/audit tables, SDKs, hosted deployment, distributed runtime, and direct Nakama/Pitaya API compatibility deferred.

## User-Facing Need

The project is aiming at a Nakama-class game/backend server framework shaped for AI-native development and testing. Friends relationships are a core social graph primitive, but adapter behavior should be explicit and testable before it is implemented.

This slice gives future agents a precise persistence-adapter contract so a later implementation can be written from a bounded spec and test plan.

## Stop Conditions

- PostgreSQL adapter implementation or SQL execution is required.
- Runtime friend request/list/status behavior is required.
- Protocol routes, Protobuf source, or generated output are required.
- Authentication/session behavior, request identity validation, or startup wiring changes are required.
- Event/audit tables, chat, groups, parties, matchmaking, match runtime, SDKs, hosted deployment, distributed runtime, or direct compatibility are required.
- Verification fails.

