# Player Account And Session Contract Standard

Status: Draft v0.2
Last updated: 2026-05-14
Scope: Player account lifecycle contracts and runtime session validation contracts
Depends on: `docs/player-identity-session-boundary.md`

The paired Simplified Chinese translation is `docs/player-account-session-contracts.zh-CN.md`. The English file is authoritative.

## 1. Purpose

This standard defines how vibit should ratify player account and session contracts before implementation.

The previous boundary standard separated player identity, player accounts, authentication, runtime session validation, transport connection metadata, envelope session metadata, and request identity context. This standard turns that boundary into a contract ratification path.

The goal is not to implement production login yet. The goal is to make the next contracts explicit enough that agents can safely add them without inventing authentication, token, persistence, or protocol behavior in the wrong layer.

Milestone status:

```text
M-005 completed on 2026-05-14.
```

The milestone completed contract ratification for the first player account and runtime session validation surfaces. It did not choose or implement authentication schemes, token behavior, credential storage, player account persistence, session persistence, Protobuf envelope changes, WebSocket handshake authentication, runtime player account handlers, WebSocket routes, direct Nakama/Pitaya API compatibility, or a major external framework dependency.

## 2. Reference Alignment

Nakama is the main reference for account, user, authentication, and session capability coverage.

Pitaya is the main reference for session binding, route handler, frontend/backend server role, and realtime server vocabulary.

vibit uses those projects as references, not governing API shapes.

Rules:

- Record each referenced pattern as adopted, adapted, deferred, or rejected.
- Do not copy Nakama or Pitaya public APIs without a compatibility ADR.
- Preserve vibit's module ownership, contract-first behavior, generated output rules, and repository checks.
- Do not let transport handlers own account, authentication, token, credential, or session persistence behavior.

## 3. Contract Families

`M-005` has two contract families.

### Player Account Contracts

Owned by:

```text
modules/player
```

Purpose:

```text
Define stable player identity and durable player account lifecycle behavior.
```

Candidate future contract vocabulary:

- `CreatePlayerAccount`
- `GetPlayerAccount`
- `LinkPlayerAccountIdentity`
- `DisablePlayerAccount`
- `DeletePlayerAccount`
- `PlayerAccountCreated`
- `PlayerAccountLinked`
- `PlayerAccountDisabled`
- `PlayerAccountDeleted`

These names are candidate vocabulary only until a later work item ratifies concrete public contracts.

### Runtime Session Validation Contracts

Owned by:

```text
runtime/internal/app
```

Purpose:

```text
Define how a decoded request becomes a validated request identity before domain handlers run.
```

Candidate future contract vocabulary:

- `ValidateSession`
- `BindSessionToConnection`
- `RefreshSession`
- `InvalidateSession`
- `SessionValidated`
- `SessionInvalidated`

These names are candidate vocabulary only until a later work item ratifies concrete public contracts.

## 4. Ownership Rules

The player module owns:

- Stable `player_id` semantics.
- Player account lifecycle.
- Player account public commands, queries, events, errors, and permissions after ratification.
- Future player account repository interfaces after persistence is ratified.

The player module does not own:

- WebSocket connections.
- Protobuf envelope framing.
- Token parsing or token signing.
- Credential storage or password hashing.
- Runtime session persistence until separately ratified.
- Inventory, currency, rewards, quests, matches, rooms, chat, or presence.

Application dispatch owns:

- Request identity context.
- Session validation handoff interfaces.
- The point where metadata-only identity can be replaced by validated identity.

Application dispatch does not own:

- Player account lifecycle.
- Authentication provider implementation.
- Credential storage.
- WebSocket connection mechanics.
- Generated Protobuf packages.

Transport owns:

- Accepted WebSocket connections.
- Binary frame IO.
- Transport-local connection metadata.

Transport does not own:

- Authentication.
- Player accounts.
- Tokens.
- Credentials.
- Domain permissions.
- Durable session state.

## 5. Reference Pattern Map

### Nakama Patterns

| Pattern | Vibit position | Reason |
| --- | --- | --- |
| User/account as first-class backend capability | Adopted | vibit needs a durable owner for stable player identity and account lifecycle. |
| Multiple authentication methods | Deferred | Supported login methods must be selected separately because they affect security, storage, and public API shape. |
| Session token and refresh token concepts | Deferred | Token format, signing, refresh, expiration, and revocation require a separate decision. |
| Realtime socket bound to authenticated session | Adapted | vibit will validate request identity before domain dispatch, but WebSocket handshake behavior remains unchanged for now. |
| Broad social, storage, leaderboard, matchmaker capability surface | Deferred | Useful roadmap input, but outside player account/session contract ratification. |

### Pitaya Patterns

| Pattern | Vibit position | Reason |
| --- | --- | --- |
| Session object separate from transport acceptor | Adopted | vibit keeps WebSocket transport metadata separate from application request identity. |
| Route handler receives session context | Adapted | vibit handlers receive `RequestIdentity` through application dispatch rather than a Pitaya API shape. |
| Frontend/backend server role split | Deferred | Distributed topology remains deferred until modular monolith boundaries are stable. |
| Groups and push vocabulary | Deferred | Useful for future room, presence, and broadcast work, but not part of account/session contract ratification. |
| Direct Pitaya API compatibility | Rejected for now | vibit defines agent-readable contracts and may add compatibility only through a future ADR. |

## 6. Ratification Order

The recommended contract ratification order is:

1. Standardize account/session contract rules and reference mapping.
2. Ratify a minimal player account semantic contract set.
3. Ratify runtime session validation semantic contracts.
4. Ratify Protobuf request/response messages for account contracts that need client/server wire shape.
5. Decide whether player account persistence and session persistence are needed.
6. Implement only after the relevant contracts and ownership rules are registered.

This order is intentionally conservative. It lets agents keep moving without silently choosing authentication or storage decisions.

The executed order may ratify player account Protobuf wire messages before runtime session semantic contracts when all of the following are true:

- The player account semantic contracts are already ratified.
- No runtime player handlers or WebSocket routes are added.
- The Protobuf envelope and WebSocket handshake remain unchanged.
- Authentication, token behavior, credential storage, session persistence, and player account persistence remain deferred.

This exception was used for the first player account wire shape so generated Protobuf output could be checked before runtime session contracts were completed.

## 7. First Minimal Contract Set

The first minimal player account contract set should prefer account lifecycle without login method selection.

Ratified first semantic contracts:

- `CreatePlayerAccount`
- `GetPlayerAccount`
- `PlayerAccountCreated`
- `player_account_errors`
- `player_account_permissions`

Recommended account fields at contract level:

- `player_id`
- `display_name`
- `account_state`
- `created_at`
- `updated_at`

The first minimal session validation contract set should prefer runtime handoff vocabulary without token format selection.

Ratified first runtime session validation semantic contracts:

- `ValidateSession`
- `SessionValidated`
- `session_errors`
- `session_permissions`

Recommended session validation fields at contract level:

- `session_id`
- `player_id`
- `connection_id`
- `validation_status`
- `actor_kind`
- `validated_at`

The player account contracts listed above have ratified semantic contracts and first Protobuf wire messages.

The runtime session validation contracts listed above are ratified semantic contracts only. They are owned by `runtime/internal/app` and sourced under:

```text
contracts/runtime/session/
```

They describe the existing application-owned session validation handoff and preserve the current metadata-only validator behavior. They do not implement real authentication, token validation, credential lookup, session persistence, player account lookup, Protobuf envelope changes, or WebSocket handshake authentication.

Ratified first player account Protobuf messages:

- `CreatePlayerAccountRequest`
- `CreatePlayerAccountResponse`
- `GetPlayerAccountRequest`
- `GetPlayerAccountResponse`
- `PlayerAccountCreated`

The ratified player account Protobuf source is:

```text
proto/vibit/player/v1/player.proto
```

The generated Go Protobuf output is:

```text
runtime/internal/generated/proto/vibit/player/v1/player.pb.go
```

Database columns, indexes, token claims, generated dispatch shapes, runtime handlers, and WebSocket handshake fields must still be ratified separately.

## 8. Required Contract Artifacts

When a player public contract is ratified, update:

- `contracts/player/...`
- `.arch/contracts.yaml`
- `modules/player/module.yaml`
- `docs/player-account-session-contracts.md`
- `docs/player-account-session-contracts.zh-CN.md`
- The relevant change spec under `changes/`

When runtime session validation public contracts are ratified, update:

- `contracts/runtime/session/...`
- `.arch/contracts.yaml`
- `.arch/runtime.yaml`
- `docs/player-account-session-contracts.md`
- `docs/player-account-session-contracts.zh-CN.md`
- The relevant change spec under `changes/`

If Protobuf messages are added, update:

- `proto/vibit/<module>/v1/...`
- `buf.yaml` or `buf.gen.yaml` only if generation roots change.
- Generated output declarations.
- Protocol alignment checks.

For the first player account wire shape, generation roots did not change. `buf generate` produced only Go Protobuf output from the new player Protobuf source.

## 9. Ask-First Boundaries

Ask the maintainer before:

- Choosing concrete login methods such as guest, device, email/password, social login, custom ID, or external identity providers.
- Choosing JWT, opaque tokens, refresh tokens, signing, expiration, or revocation behavior.
- Adding credential storage, password hashing, cryptography, OAuth, OIDC, or external auth dependencies.
- Adding player account database schema, migrations, or session persistence.
- Changing the Protobuf envelope.
- Changing WebSocket handshake authentication behavior.
- Making metadata-only `player_id` or `session_id` sufficient for production permission grants.
- Copying Nakama or Pitaya public API shape.

## 10. Verification

Before implementing any contract from this standard, run:

```bash
node tools/vibit check contracts --json
node tools/vibit check protocol --json
node tools/vibit check runtime --json
node tools/vibit check work --json
```

After adding public contracts, run:

```bash
node tools/vibit check all --json
```

If Go runtime code changes, also run:

```bash
cd runtime && go test ./...
```

## 11. Next Direction Gate

After `M-005`, agents must stop at `M-006/W-0037` before choosing the next major implementation direction.

The next direction may be one of the following, but this standard does not choose among them:

- Authentication and token/session validation design.
- Player account PostgreSQL schema and persistence.
- Runtime player account handlers and WebSocket route wiring.
- Additional core game backend modules, reviewed against Nakama and Pitaya capability coverage.
- Generator and contract tooling expansion before additional runtime features.

Choosing among these paths is a maintainer decision because it affects long-lived security, data ownership, protocol, runtime, or roadmap shape.
