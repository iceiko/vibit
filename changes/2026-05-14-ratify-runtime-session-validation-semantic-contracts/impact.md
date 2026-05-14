# Impact

## Affected Areas

- `runtime/internal/app`: owns the session validation handoff boundary described by these semantic contracts.
- `contracts/runtime/session`: receives runtime-owned contract source files.
- `.arch/contracts.yaml`: registers runtime session contracts separately from domain module contracts.
- `.arch/runtime.yaml`: records the ratified semantic contract state.
- `.arch/work-items.yaml`: marks `W-0034` complete and exposes the next bounded step if safe.

## Module Ownership Impact

No domain module ownership moves.

The player module still owns player identity and player account lifecycle. Runtime application dispatch owns request identity context and the session validation handoff. WebSocket transport remains outside authentication and session validation.

## Public Contract Impact

Adds semantic runtime session contracts:

- `ValidateSession`
- `SessionValidated`
- `session_errors`
- `session_permissions`

These contracts are runtime/application boundary contracts, not player module contracts.

## Data And Migration Impact

No database schema, migration, repository, credential store, token store, or session store is added.

## Protocol Impact

No Protobuf source is added. The existing envelope and WebSocket handshake remain unchanged.

## Runtime Impact

No Go implementation changes are required. The existing `runtime/internal/app/SessionValidator` hook remains the runtime extension point. Metadata-only validation remains current behavior and is explicitly not authenticated proof.

## Compatibility Risks

Low. This change makes existing runtime session validation semantics explicit before implementation. The main risk is future agents mistaking the semantic contract for production authentication; the contract and checks explicitly block that.
