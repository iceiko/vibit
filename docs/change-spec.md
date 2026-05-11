# Change Spec Standard

Status: Draft v0.1  
Last updated: 2026-05-12  
Scope: `changes/<date>-<change-id>/`

This document defines the standard change-spec workflow for vibit.

A change spec is durable working context for humans and agents. It exists so non-trivial changes are not performed from a short prompt alone.

## 1. Purpose

Every non-trivial feature, bug fix, migration, refactor, or standard change should have a directory:

```text
changes/<date>-<change-id>/
```

Example:

```text
changes/2026-05-12-add-inventory-module/
```

The directory should preserve:

- The original request
- The clarified requirement
- Affected modules and contracts
- Impact analysis
- Implementation plan
- Acceptance checklist
- Verification results
- Open questions

## 2. Required Files

Recommended structure:

```text
changes/<date>-<change-id>/
  request.md
  spec.yaml
  impact.md
  plan.md
  checklist.md
  verification.md
```

For small documentation-only changes, a lighter version may be acceptable, but agents must still record what changed and how it was verified.

## 3. `request.md`

`request.md` captures the human request and the clarified requirement.

It should include:

- Original request
- Clarified requirement
- User-visible outcome
- Non-goals
- Unknowns
- Acceptance criteria

## 4. `spec.yaml`

`spec.yaml` is the machine-readable summary of the change.

Example:

```yaml
schema_version: 0.1
change_id: add-inventory-module
date: 2026-05-12
type: feature
status: draft

summary: Add the first inventory module prototype.

affected_modules:
  - inventory

contracts:
  commands:
    added:
      - AddItem
      - RemoveItem
    changed: []
    removed: []
  queries:
    added:
      - GetInventory
    changed: []
    removed: []
  events:
    added:
      - ItemAdded
      - ItemRemoved
    changed: []
    removed: []
  permissions:
    added:
      - inventory_read
      - inventory_write
    changed: []
    removed: []

data:
  migrations_required: false
  ownership_changes: []

compatibility:
  breaking_api: false
  breaking_events: false
  breaking_data: false

verification:
  required:
    - architecture_checks
    - command_tests
    - invariant_tests
  status: Not applicable
```

## 5. `impact.md`

`impact.md` explains what the change touches and why.

It should cover:

- Affected modules
- Module ownership impact
- Public contract impact
- Event impact
- Permission impact
- Data and migration impact
- Test impact
- Documentation impact
- Compatibility risks

## 6. `plan.md`

`plan.md` is the bounded implementation plan.

It should identify:

- Files to create
- Files to edit
- Generated artifacts
- Handwritten logic
- Tests to add or update
- Verification commands
- Rollback or migration notes when relevant

## 7. `checklist.md`

`checklist.md` tracks completion.

Use simple task states:

```text
- [ ] Pending task
- [x] Completed task
```

The checklist should include contract, implementation, test, verification, and documentation tasks.

## 8. `verification.md`

`verification.md` records what was checked.

Use this format:

```text
Verified:
- <command or check>

Not verified:
- <reason>

Not applicable:
- <reason>
```

Never claim a change is verified when verification did not run.

## 9. Agent Rules

Agents should create or update a change spec when:

- The change affects public behavior.
- The change affects module ownership.
- The change adds or changes commands, queries, events, permissions, errors, or data shape.
- The change introduces or changes an architectural standard.
- The change is large enough that another agent would need durable context to continue.

Agents may skip a full change spec for small typo fixes, formatting-only edits, or narrow documentation updates, but the final response must still describe verification status.

## 10. Naming Rules

Change directory names should use:

```text
YYYY-MM-DD-short-kebab-case-id
```

Rules:

- Use the date when the change spec is created.
- Use a short, stable, descriptive ID.
- Do not rename after implementation starts unless necessary.
- Keep the ID meaningful to future agents.
