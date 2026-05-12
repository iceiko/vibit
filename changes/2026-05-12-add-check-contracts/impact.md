# Impact

## Affected Modules

`inventory` is affected because it currently owns the only registered contract source files.

## Module Ownership Impact

No ownership changes are made.

## Public Contract Impact

No public runtime contracts are added or changed.

The change adds tooling that checks existing public contract source files.

## Event Impact

No event contracts are changed.

## Permission Impact

No permission contracts are changed.

## Data And Migration Impact

No data or migration impact.

## Test Impact

The CLI check itself is verified by running:

- `node tools/vibit check contracts`
- `node tools/vibit check contracts --json`
- `node tools/vibit check all --json`

No runtime tests are added.

## Documentation Impact

README and AGENTS command lists are updated in English and Simplified Chinese.

## Compatibility Risks

Low.

The check is additive. It may make future inconsistent contract edits fail earlier, which is intended.
