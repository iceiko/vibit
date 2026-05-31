# Request

Define the friends relationship runtime behavior gate for the Nakama-first friends/groups/parties capability.

The gate must stay inside the runtime-behavior-boundary slice:

- define the future application-owned runtime behavior owner, package candidate, service source candidate, and test candidate;
- preserve the storage-neutral `runtime/internal/modules/friends.Repository` interface and implemented PostgreSQL adapter as lower-level dependencies;
- define validated request identity derivation for actor identity;
- define actor-relative public status behavior for request, accept, reject, remove, block, unblock, list, and status operations;
- define permission and route-policy posture, repository/unit-of-work handoff, conflict mapping, redaction, and test expectations;
- open a bounded follow-up for implementation;
- keep actual runtime behavior implementation, handlers, startup wiring, protocol routes, Protobuf source, generated output, dependencies, migration changes, event/audit tables, SDKs, hosted deployment, distributed runtime, and direct Nakama/Pitaya API compatibility deferred.

## User-Facing Need

The project is aiming at a Nakama-class game/backend server framework shaped for AI-native development and testing. Friends relationships are a core social graph primitive, but actor identity, privacy, route policy, and conflict behavior should be explicit and testable before runtime behavior is implemented.

This slice gives future agents a precise application-behavior contract so a later implementation can be written from a bounded spec and test plan.

## Stop Conditions

- Runtime friendship behavior implementation is required.
- Runtime handlers, command/query registration, or startup composition are required.
- Protocol routes, Protobuf source, or generated output are required.
- Authentication/session behavior, request identity validation, or route-policy semantics changes are required.
- Repository interface, PostgreSQL adapter, or migration changes are required.
- Event/audit tables, chat, groups, parties, matchmaking, match runtime, SDKs, hosted deployment, distributed runtime, or direct compatibility are required.
- Verification fails.

