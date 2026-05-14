# Verification

Verified:

- `node tools/vibit inspect next --json`
  - Passed and reported `W-0089` as the next ready work item.
- `node tools/vibit check work --json`
  - Passed: 545 checks
- `node tools/vibit check change decide-authentication-generated-contract-shape-timing --json`
  - Passed: 13 checks
- `node tools/vibit check generated --json`
  - Passed: 97 checks
- `node tools/vibit check contracts --json`
  - Passed: 290 checks
- `node tools/vibit inspect generated --json`
  - Reported 3 generated Protobuf files, 10 generated contract shape files, 10 expected existing contract shapes, 0 missing shapes, and 0 stale files.
- `node tools/vibit check runtime --json`
  - Passed: 662 checks
  - Warning: 1 existing `runtime.identity_boundary` warning for authentication PostgreSQL adapter persistence vocabulary inside the ratified boundary.
- `node tools/vibit check module authentication --json`
  - Passed: 23 checks
- `node tools/vibit check memory --json`
  - Passed: 1146 checks
- `node tools/vibit check all --json`
  - Passed: 140 subchecks
  - Warning: 1 existing runtime identity-boundary warning noted above
- `git diff --check`
- `rg -n -l "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" .git/config . --glob '!.git/**' --glob '!.vibit.local.env'`
  - No matches

Not verified:

- None.

Not applicable:

- Runtime Go tests are not required because this change does not add or modify Go runtime behavior.
