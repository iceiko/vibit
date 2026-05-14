# Player Identity And Session Boundary Standard

Status: Draft v0.1
Last updated: 2026-05-14
Scope: Player identity, account, authentication, WebSocket session, and request identity context boundaries
Canonical decision: `ADR-0021`

The paired Simplified Chinese translation is `docs/player-identity-session-boundary.zh-CN.md`. The English file is authoritative.

## 1. Purpose

vibit's first durable inventory runtime can already process player-scoped commands and queries, but the current `player_id` is still request data rather than authenticated runtime identity.

This standard defines the boundary before more modules depend on ad hoc player context. It prevents future agents from placing identity, authentication, session validation, or player-account behavior in whichever layer is easiest to edit.

The goal of this step is boundary clarity, not a production login system.

## 2. External Reference Alignment

Nakama is the main reference for the product capability surface around accounts, authentication, users, and sessions.

Pitaya is the main reference for Go game server vocabulary around acceptors, session binding, route handlers, and server push.

vibit adopts the separation of concerns, not their public APIs:

- Account and authentication concepts must not be hidden in transport handlers.
- Connection sessions and route dispatch must remain separate concerns.
- Realtime connection metadata is not the same thing as durable player account ownership.
- Any future API compatibility with Nakama, Pitaya, or another project requires a separate ADR.

## 3. Identity Vocabulary

The following terms are distinct.

### Player Identity

`player_id` is the stable domain identity used by game modules to address a player-owned aggregate.

Owner:

```text
player module
```

Rules:

- `player_id` is a domain identity, not a WebSocket connection id.
- Domain modules such as inventory may reference `player_id`.
- Domain modules other than `player` must not create, authenticate, merge, delete, or rename players.
- Until the player module exists, `player_id` remains an external reference with an explicit boundary note.

### Player Account

A player account is the durable account record and lifecycle for a player.

Owner:

```text
player module
```

Rules:

- Account creation, lookup, linking, disabling, deletion, and recovery belong to the player module.
- Account data must not be stored in inventory tables or transport/session structures.
- Adding player account database migrations requires a separate change spec and maintainer confirmation.

### Authentication

Authentication proves that an actor may bind to a player identity or privileged service identity.

Owner:

```text
authentication boundary, implemented through the player/session capability when ratified
```

Rules:

- This standard does not choose a concrete authentication scheme.
- Do not introduce JWT, OAuth, OIDC, password hashing, password storage, guest login, device login, social login, or external identity providers without a separate decision.
- Authentication must produce a machine-readable result for downstream session binding instead of being interpreted inside domain handlers.

### Runtime Session

A runtime session is the server-side logical session associated with an accepted client connection after session validation exists.

Owner:

```text
runtime session boundary
```

Rules:

- A runtime session may bind a connection to `session_id`, `player_id`, and future authorization claims.
- A runtime session is not a player account.
- A runtime session is not owned by the WebSocket transport adapter.
- Session storage, expiration, reconnect behavior, and token validation are deferred until separate ratification.

### Transport Connection

`connection_id` is transport-local connection metadata.

Owner:

```text
runtime/internal/platform/transport/ws
```

Rules:

- `connection_id` may change across reconnects.
- Transport may expose connection metadata, but it must not authenticate it as a player.
- Transport must not parse credentials, domain payloads, or player account data.

### Request Identity Context

Request identity context is the application-facing identity data attached to a dispatched command or query.

Owner:

```text
runtime/internal/app
```

Rules:

- Request identity context is built after protocol decoding and session validation.
- Domain handlers receive already-normalized identity context through vibit-owned application request types.
- A module may compare request identity to requested target data when enforcing permissions.
- A module must not validate transport credentials or token formats.

## 4. Layer Ownership

### WebSocket Transport

Owner:

```text
runtime/internal/platform/transport/ws
```

Responsibilities:

- Accept WebSocket connections.
- Read and write binary frames.
- Provide transport-local connection metadata.
- Delegate opaque frame bytes to injected protocol/application composition.

Must not:

- Authenticate players.
- Parse tokens, credentials, or account payloads.
- Create or validate player accounts.
- Treat `connection_id` as durable identity.
- Enforce domain permissions.

### Protobuf Protocol Adapter

Owner:

```text
runtime/internal/platform/protocol/protobuf
```

Responsibilities:

- Decode and encode the existing Protobuf envelope.
- Preserve session metadata fields that already exist in the envelope.
- Convert envelope metadata into application handoff types.
- Map malformed identity/session metadata into protocol or application errors when validation rules exist.

Must not:

- Choose or implement authentication schemes.
- Own long-lived session state.
- Change the envelope shape without a protocol change spec and ADR.
- Enforce module-level permissions.

### Application Dispatch And Session Validation

Owner:

```text
runtime/internal/app
```

Responsibilities:

- Receive decoded route requests.
- Invoke session validation once the session boundary is implemented.
- Pass normalized request identity context to command and query handlers.
- Keep route dispatch separate from authentication provider details.

Must not:

- Store player account lifecycle data.
- Import WebSocket libraries.
- Import generated Protobuf packages.
- Hide module business rules.

### Player Module

Future owner:

```text
modules/player
runtime/internal/modules/player
```

Responsibilities:

- Own player identity and player account lifecycle.
- Define public account/session commands, queries, events, permissions, and errors before implementation.
- Own persistent player account state when migrations are ratified.
- Publish account lifecycle events when needed by other modules.

Must not:

- Own inventory state.
- Own WebSocket connection mechanics.
- Hide transport or protocol behavior.

### Inventory Module

Owner:

```text
modules/inventory
runtime/internal/modules/inventory
```

Responsibilities:

- Own inventory records and item quantities keyed by `player_id`.
- Enforce inventory permissions and invariants.
- Treat `player_id` as an external reference until the player module exists.

Must not:

- Create, authenticate, link, disable, delete, or migrate player accounts.
- Own session state or token validation.
- Depend directly on the player module until a future change explicitly allows the dependency.

## 5. Current Envelope Position

The existing protocol envelope already contains:

```text
Session.connection_id
Session.session_id
Session.player_id
Session.connection_epoch
```

This standard does not change those fields.

Current meaning:

- `connection_id`: transport-local metadata when available.
- `session_id`: reserved authenticated logical session identifier.
- `player_id`: reserved authenticated player identity when available.
- `connection_epoch`: reserved reconnect/lifecycle version.

Until session validation exists, these fields are envelope metadata. They must not be treated as proof of identity by domain logic.

## 6. Request Processing Model

The intended future request flow is:

```text
websocket frame
-> transport connection metadata
-> protobuf envelope
-> preliminary route request
-> session validation
-> request identity context
-> application dispatch
-> domain command or query handler
```

Current runtime state:

- WebSocket transport and Protobuf request dispatch exist.
- Inventory handlers can receive `app.Session` and `app.Target` metadata.
- An application-owned session validation hook exists with a metadata-only default path.
- Inventory permission policies receive an application-owned request identity handoff context, but current bootstrap policy remains explicitly static.
- No real authentication, token parsing, credential lookup, session persistence, or player account lookup exists yet.
- Metadata-only request identity is not authenticated proof and must not satisfy identity-aware privileged permission policy.

## 7. Permission Boundary

Permissions are module-owned business checks.

Authentication and session validation answer:

```text
Who is the actor?
```

Module permission policies answer:

```text
May this actor perform this operation on this module target?
```

Inventory must move away from bootstrap static allow behavior through a future work item, but that future item must still avoid choosing an authentication scheme unless separately ratified.

For inventory:

- `GrantItem` requires `inventory_grant_item`.
- `GetInventory` requires `inventory_read`.
- Inventory permission context carries the requested actor text, target `player_id`, and `RequestIdentity`.
- Static bootstrap permission policy may explicitly allow current local proof-slice behavior, but it is not production authorization.
- Identity-aware permission policies must treat `metadata_only` identity as unauthenticated.
- A future player-bound read policy may allow a player to read only their own inventory.
- Service or admin actors require explicit actor/permission modeling before privileged grants are production behavior.

## 8. Deferred Decisions

The following are intentionally not decided by this standard:

- Token format.
- JWT, OAuth, OIDC, password, device, guest, social, or external-provider authentication.
- Credential storage.
- Session persistence store.
- Session expiration and refresh model.
- Reconnect replay behavior.
- Player account database schema.
- Player account migration files.
- Protobuf envelope changes.
- WebSocket handshake authentication contract.
- Presence, parties, rooms, matches, or broadcast group behavior.

Each requires a separate change spec and, where architectural, an ADR.

## 9. Application Handoff

The current Go runtime defines an application-owned `RequestIdentity` handoff under:

```text
runtime/internal/app
```

Rules:

- `RouteRequest.Identity` is the application-facing identity context for command and query handlers.
- `ApplicationResult.Identity` preserves the request identity context for downstream protocol mapping and future auditing.
- `MetadataOnlyIdentityFromSession` may normalize existing session metadata into request identity context.
- `MetadataOnlySessionValidator` preserves current metadata-only behavior without authenticating clients.
- `SessionValidatingDispatcher` is the application-owned hook that runs after protocol decoding and before module handlers receive the request.
- Metadata-only identity is not authenticated identity.
- `PlayerIDValidated` and `SessionValidated` remain false until a future session validator performs real validation.
- `SessionValidationResult` is handoff vocabulary and hook output. It is not an authentication implementation, token contract, or session store.
- Future real session validation should be implemented behind this hook, replacing metadata-only identity with validated identity when validation succeeds.
- Inventory consumes `RequestIdentity` only through its permission handoff context. Inventory must not validate tokens, query player accounts, own sessions, or import the player module.

## 10. Completed Boundary Sequence

The M-003 boundary sequence is complete:

1. Player identity and session responsibilities are separated in this standard and `ADR-0021`.
2. `modules/player/module.yaml` and module agent guides declare player identity and account lifecycle ownership without runtime code, migrations, credential storage, or authentication providers.
3. Application-owned request identity and session validation handoff types exist under `runtime/internal/app`.
4. `MetadataOnlySessionValidator` and `SessionValidatingDispatcher` provide the future validation hook while preserving metadata-only behavior.
5. Inventory permission policies receive `RequestIdentity` through `PermissionContext`, and `MetadataOnlyDenyPermissionPolicy` prevents metadata-only identity from satisfying privileged grants.
6. `runtime.identity_boundary` repository checks protect the boundary from common transport, domain, generated Protobuf, auth dependency, player runtime, unauthorized player Protobuf, and player persistence regressions.

Next implementation must not proceed until the maintainer confirms which major direction comes next.

Plausible next milestones include:

- Ratify player account and session public contracts.
- Ratify authentication, token, and session validation design.
- Continue game-domain breadth with item catalog, currency, rewards, quests, or match sessions before production authentication.
- Improve generators and contract tooling before expanding runtime features.

Agents must stop before any step that chooses a concrete authentication mechanism, token format, credential store, player account schema, session persistence model, Protobuf envelope change, or WebSocket handshake contract.

## 11. Verification

Current repository checks:

```bash
node tools/vibit check architecture
node tools/vibit check protocol
node tools/vibit check runtime
node tools/vibit check work
node tools/vibit check all
```

`node tools/vibit check runtime` includes the `runtime.identity_boundary` repository check.

The check currently verifies:

- WebSocket transport does not import runtime domain modules, player runtime packages, inventory runtime packages, generated Protobuf packages, or Protobuf runtime dependencies.
- Domain modules do not import WebSocket transport, generated Protobuf packages, Protobuf runtime dependencies, or known authentication, token, OAuth, OIDC, credential, and password-hashing dependencies.
- Runtime player module code remains absent until public player contracts are ratified.
- Player Protobuf source roots are allowed only for ratified player wire contracts and must not imply runtime handlers, authentication, or persistence.
- PostgreSQL migrations do not introduce player identity or player account persistence before schema ratification.
- `modules/player/module.yaml` preserves identity/session boundary markers after semantic and wire contracts are ratified.

Intentionally deferred checks:

- Proving semantic authentication correctness. The current repository has no authentication implementation to verify.
- Proving all possible third-party authentication libraries. The current check blocks known dependency families and must be expanded when a concrete authentication decision is ratified.
- Inspecting generated player/session output beyond the first ratified player account Protobuf package. Runtime session generated output does not exist yet.
- Proving permission policy completeness for production authorization. Current policies only define the handoff and metadata-only guard boundary.

## 12. Agent Rules

Agents must:

- Read this standard before adding player, account, authentication, session, permission, or request identity behavior.
- Keep player identity, account lifecycle, authentication, runtime session, transport connection, and request identity context separate.
- Add or update manifests and contracts before implementing public behavior.
- Record any adopted, adapted, deferred, or rejected Nakama/Pitaya pattern in the relevant change spec or ADR.
- Ask before crossing a deferred decision boundary.

Agents must not:

- Treat `player_id` supplied by a client envelope as authenticated proof before session validation exists.
- Add player account migrations during boundary-only work.
- Add JWT, OAuth, OIDC, password hashing, password storage, or external identity provider dependencies without ratification.
- Move inventory ownership into the player module.
- Hide authentication or permission checks inside WebSocket transport or Protobuf bridge code.
