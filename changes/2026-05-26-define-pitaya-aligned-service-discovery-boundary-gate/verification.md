# Verification

## RED Checks

The initial checks failed because the W-0252 rule and change artifacts were absent:

```text
node tools/vibit inspect rule runtime.pitaya_aligned_service_discovery_boundary_gate
# Unknown rule_id: runtime.pitaya_aligned_service_discovery_boundary_gate

node tools/vibit check change define-pitaya-aligned-service-discovery-boundary-gate --json
# change directory does not exist
```

## Required Verification

- `node -c tools/vibit`
- `node tools/vibit inspect next --json`
- `node tools/vibit inspect rule runtime.pitaya_aligned_service_discovery_boundary_gate`
- `node tools/vibit check change define-pitaya-aligned-service-discovery-boundary-gate --json`
- `node tools/vibit check work --json`
- `node tools/vibit check runtime --json`
- `node tools/vibit check memory --json`
- `node tools/vibit check schemas --json`
- `node tools/vibit check all --json`
- `git diff --check`

## Result

Passed on 2026-05-31.

```text
node -c tools/vibit
# exit 0

node tools/vibit inspect next --json
# status: ready
# next_ready: W-0253 Implement Pitaya-aligned service discovery source-first map

node tools/vibit inspect rule runtime.pitaya_aligned_service_discovery_boundary_gate
# rule found in rules/check-rules.json

node tools/vibit check change define-pitaya-aligned-service-discovery-boundary-gate --json
# passed: 13, warnings: 0, failures: 0

node tools/vibit check work --json
# passed: 1530, warnings: 0, failures: 0

node tools/vibit check runtime --json
# passed: 20681, warnings: 1, failures: 0
# known warning: runtime.identity_boundary in runtime/internal/platform/persistence/postgres/authentication_repository.go

node tools/vibit check memory --json
# passed: 4292, warnings: 0, failures: 0

node tools/vibit check schemas --json
# passed: 4700, warnings: 0, failures: 0

node tools/vibit check all --json
# subchecks: 307, passed: 307, warnings: 1, failures: 0

git diff --check
# exit 0
```
