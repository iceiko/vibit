# Verification

Verified:

- `node tools/vibit --help`
  - Result: passed; help includes `check all`.
- `node tools/vibit check all`
  - Result: passed; ran architecture, all discovered change specs, and registered module checks.
- `node tools/vibit check architecture`
  - Result: passed.
- `node tools/vibit check change add-vibit-check-all`
  - Result: passed.
- `node tools/vibit check module inventory`
  - Result: passed.
- `rg -n "ghp_[A-Za-z0-9]|github_pat_[A-Za-z0-9]" .`
  - Result: no raw GitHub token pattern was found.

Not verified:

- Failure-path behavior for `check all` was not exhaustively tested with intentionally corrupted files.

Not applicable:

- Runtime server tests do not apply because this change targets tooling.
