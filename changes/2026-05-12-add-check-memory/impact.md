# Impact Analysis

## Affected Modules

No domain modules are affected.

## Module Ownership Impact

No ownership changes.

## Public Contract Impact

Adds a new CLI check command:

```bash
node tools/vibit check memory
node tools/vibit check memory --json
```

Changes `check all` by adding `check memory` as an aggregate subcheck.

Adds rule catalog entries:

- `artifact.minimum_count`
- `markdown.heading`
- `markdown.metadata`

## Event Impact

No events are added, changed, or removed.

## Permission Impact

No permissions are added, changed, or removed.

## Data And Migration Impact

No runtime data or migrations.

## Test Impact

Verification must confirm text and JSON output for `check memory`, and confirm `check all` still passes.

## Documentation Impact

Update:

- README
- AGENTS
- Rule catalog
- Change spec
- Conversation log

## Compatibility Risks

Low. This is additive, but `check all` now enforces memory artifacts. Existing repository memory must satisfy the new rule before this change can be accepted.
