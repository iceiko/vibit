# Impact

## Scope

This change adds a source-first repository inspection map for Pitaya-aligned cluster-safe session routing vocabulary.

Affected areas:

- `tools/vibit`
- `rules/check-rules.json`
- architecture manifests under `.arch/`
- repository memory docs and module guidance
- W-0257 change artifacts and ADR-0165

## Behavior

No runtime behavior changes are added.

No protocol shape changes are added.

No generated output is added.

No dependencies are added.

No persistence, repository interface, PostgreSQL adapter, or migration behavior is added.

## Risk

The main risk is semantic drift: future agents could interpret session-routing vocabulary as permission to implement cluster routing. The inspection output, rule text, manifests, ADR, and deferral flags explicitly prevent that.

## Follow-Up

W-0258 selects the next Pitaya-aligned direction after this source-first map.
