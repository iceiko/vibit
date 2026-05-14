# Impact

## Affected Modules

- `inventory`
- `player`

Only generated contract shape metadata is added for these modules.

## Module Ownership Impact

No ownership changes.

## Public Contract Impact

No semantic contracts are added or changed. Existing registered contracts are summarized into generated Go shape files.

## Runtime Impact

No runtime behavior is added. The generated files are not handler implementations.

## Generated Output Impact

Adds generated files under:

- `runtime/internal/generated/contracts/inventory/`
- `runtime/internal/generated/contracts/player/`

## Compatibility Risks

Low. Generated files are additive and metadata-only.
