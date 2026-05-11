# Impact Analysis

## Affected Modules

- None. No runtime modules exist yet.

## Module Ownership Impact

No ownership changes.

## Public Contract Impact

Adds initial CLI command surface:

- `vibit check architecture`
- `vibit check module <module>`
- `vibit check change <change-id>`
- `vibit generate module <module>`

## Data And Migration Impact

No runtime data or migrations.

## Test Impact

The first CLI should have enough verification to prove:

- Help output works.
- Architecture check can pass on the current repository.
- Change check can pass on existing change specs.
- Module check reports missing modules in a deterministic way.

## Documentation Impact

Expected documentation updates:

- README usage section
- AGENTS verification command section
- Possibly `.arch/conventions.yaml` command registry

## Compatibility Risks

Low. This is the first CLI surface. The main risk is prematurely locking in poor command behavior.
