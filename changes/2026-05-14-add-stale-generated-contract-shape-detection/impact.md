# Impact Analysis

## Affected Modules

- None directly.

## Module Ownership Impact

No module ownership changes.

## Public Contract Impact

No public contract changes.

## Data And Migration Impact

No data or migration impact.

## Test Impact

The existing generated output check now has stricter coverage for stale files.

## Documentation Impact

No public documentation change is required because stale-file detection fits the existing generated output standard.

## Compatibility Risks

Low. Repositories with stale generated contract shape files will now fail `check generated`, which is the intended guard.
