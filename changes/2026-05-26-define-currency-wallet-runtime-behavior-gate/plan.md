# Plan: Define Currency Wallet Runtime Behavior Gate

## Steps

1. Read current currency wallet repository and PostgreSQL adapter boundaries.
2. Model the gate after storage objects and friends relationship runtime behavior gates.
3. Add English and Simplified Chinese gate standards.
4. Add ADR-0207 and conversation memory.
5. Update manifests, module guide, README, alpha docs, reference roadmap, rules, and `tools/vibit`.
6. Verify the rule and full repository checks.

## Boundary

This plan intentionally excludes runtime implementation, protocol shape, generated output, dependencies, migrations, reward/inventory/purchase integration, catalog/event tables, payment behavior, hosted surfaces, SDKs, release artifacts, distributed runtime, and direct compatibility.
