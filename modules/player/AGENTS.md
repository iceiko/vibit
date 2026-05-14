# player Module Agent Guide

Status: Draft v0.1

## When To Use This Module

Use this module for requirements that define stable player identity or player account lifecycle ownership.

This module currently owns semantic and wire contract artifacts, the ratified PostgreSQL account lifecycle schema boundary, the storage-neutral runtime repository interface, and the implemented PostgreSQL adapter boundary. It owns the vocabulary and ownership boundary for:

- `player_id` as stable domain identity.
- Player account lifecycle: creation, lookup, linking, disabling, deletion, and recovery.
- Ratified player account semantic contracts.
- Ratified player account Protobuf wire messages.
- Ratified player account PostgreSQL lifecycle schema boundaries for `player_accounts` and `player_account_events`.
- The module-owned repository interface at `runtime/internal/modules/player/repository.go`.
- The PostgreSQL adapter at `runtime/internal/platform/persistence/postgres/player_account_repository.go`.

The current state includes the first player account migration source, the module-owned repository interface boundary, and the focused PostgreSQL adapter implementation. It does not implement authentication, credentials, runtime sessions, WebSocket routes, WebSocket handshake behavior, or runtime account/session handlers.

## When Not To Use This Module

Do not use this module for:

- WebSocket connection mechanics.
- Protobuf envelope routing or payload conversion.
- Runtime session validation or session storage.
- Token formats, credential parsing, password storage, JWT, OAuth, OIDC, guest login, device login, social login, or external identity providers.
- Inventory state, item quantities, item grants, or inventory persistence.
- Currency, reward, quest, match, or realtime room behavior.
- Authentication, credential, token, external identity, or session tables.
- Additional player account repository adapter behavior beyond the ratified PostgreSQL adapter scope.

If a requirement needs one of those concepts, update the owning boundary or create a separate ratified change before adding code here.

## Extension Points

- Future command handlers after public commands are ratified.
- Future query handlers after public queries are ratified.
- Future lifecycle events after public events are ratified.
- Player identity and account lifecycle policies.
- Player account repository interface: `runtime/internal/modules/player/repository.go`.
- Existing first player account migration source: `runtime/migrations/postgres/000002_create_player_account_state.sql`.
- Boundary and invariant tests.

The vocabulary placeholders in `module.yaml` are not public contracts. Do not treat them as permission to implement APIs.

Before adding any player public contract, create or update:

- `contracts/player/...`
- `.arch/contracts.yaml`
- `modules/player/module.yaml`
- A change spec under `changes/`
- Any required Protobuf source only after the protocol impact is explicit

The first ratified player account Protobuf source is `proto/vibit/player/v1/player.proto`. The generated Go Protobuf output is `runtime/internal/generated/proto/vibit/player/v1/player.pb.go`. Do not treat these wire shapes as permission to add runtime handlers, WebSocket routes, authentication, token behavior, credential storage, player persistence implementation, or session persistence.

The ratified player account PostgreSQL schema boundary is in `docs/postgresql-persistence-boundary.md` and `ADR-0022`. The first allowed player account lifecycle tables are:

- `player_accounts`
- `player_account_events`

The first player account schema must not include credentials, password hashes, authentication provider subjects, external identity links, access tokens, refresh tokens, runtime session rows, WebSocket connection state, request identity validation results, inventory state, or permission grants.

The ratified runtime repository interface is `runtime/internal/modules/player/repository.go`. It may define storage-neutral account lifecycle structs, `Repository.CreatePlayerAccount`, `Repository.GetPlayerAccount`, and the mutation/query shapes needed by persistence adapters. It must not import PostgreSQL, WebSocket, Protobuf, authentication, token, credential, session, S3, or MinIO dependencies.

The PostgreSQL adapter is implemented at `runtime/internal/platform/persistence/postgres/player_account_repository.go`, with tests at `runtime/internal/platform/persistence/postgres/player_account_repository_test.go`. It is constructed with `NewPlayerAccountRepositoryForUnitOfWork(executor)`, implements `player.Repository`, uses a caller-supplied unit-of-work executor, and avoids `BEGIN`, `COMMIT`, and `ROLLBACK`. Its first SQL scope is limited to inserting `player_accounts`, inserting `player_account_events` for `PlayerAccountCreated`, and reading current account lifecycle rows from `player_accounts`. It does not authorize runtime player handlers, WebSocket routes, authentication, tokens, credentials, external identity links, or session persistence.

## Forbidden Shortcuts

- Do not bypass boundaries declared in `module.yaml`.
- Do not directly modify data owned by another module.
- Do not add unregistered public commands, queries, events, or permissions.
- Do not make transport connection metadata act as player identity.
- Do not treat client-supplied `player_id` or `session_id` as authenticated proof.
- Do not add authentication provider, token, credential, password, JWT, OAuth, OIDC, guest login, device login, social login, or identity-provider code without a separate decision.
- Do not add another player account migration source unless a future work item specifically ratifies the schema change.
- Do not add player account repository adapters outside `runtime/internal/platform/persistence/postgres/`.
- Do not expand the player account PostgreSQL adapter beyond its ratified SQL scope without a separate work item.
- Do not let the player account PostgreSQL adapter open hidden transactions, decode Protobuf, enforce WebSocket behavior, parse credentials, validate tokens, bind sessions, or enforce permissions.
- Do not add runtime player account handlers or WebSocket routes as part of persistence work.
- Do not make this module own inventory state or depend on inventory internals.
- Do not change the Protobuf envelope or WebSocket handshake contract from this module.
- Do not import WebSocket, Protobuf, PostgreSQL, S3, or MinIO dependencies directly into the domain module.
- Do not hand-edit generated files. If generated output is wrong, change the source schema, template, or generator.

## Required Tests

See `tests.required` in `module.yaml`.

For the current semantic, wire, migration, repository-interface, and PostgreSQL-adapter-implemented state, verification is repository-level architecture, module, contracts, protocol, generated-output, runtime, work queue, change-spec checks, and focused PostgreSQL adapter behavior tests. Runtime behavior tests for player handlers become required only when Go player handlers or runtime protocol adapters are added.
