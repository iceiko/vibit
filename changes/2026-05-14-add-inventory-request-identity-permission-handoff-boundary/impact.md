# Impact Analysis

## Affected Modules

- `inventory`
- `runtime`

## Module Ownership Impact

Inventory continues to own inventory state, item grants, reads, capacity rules, and permission checks. Player identity, account lifecycle, authentication, runtime sessions, token formats, credentials, and session validation remain outside inventory.

## Public Contract Impact

No public command, query, event, Protobuf, or error shape changes are introduced. The runtime permission interface changes internally so permission policies receive a structured context that includes request identity.

## Data And Migration Impact

No data ownership changes and no migrations.

## Test Impact

Focused Go tests cover:

- Request identity reaching inventory permission policies.
- Metadata-only identity not satisfying identity-aware privileged grant policy.
- Existing bootstrap static allow behavior remaining explicit.

## Documentation Impact

Update the player/session boundary standard and inventory module guidance in English and Simplified Chinese.

## Compatibility Risks

This changes an internal Go interface. Existing runtime code in this repository will be updated in the same change. No external public API compatibility change is intended.
