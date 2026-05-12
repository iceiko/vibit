# Verification

Verified:

- `node tools/vibit inspect rules`
- `node tools/vibit inspect rules --category check`
- `node tools/vibit inspect rules --category check | node -e '<filter assertion>'`
- `node tools/vibit check all`
- `node tools/vibit check all --json`
- `rg -n "ghp_[A-Za-z0-9]|github_pat_[A-Za-z0-9]" . .git/config`
- `git diff --check`

Not verified:

- Full JSON Schema validation is not implemented yet; current schema checks validate critical fields only.

Not applicable:

- Runtime server tests do not apply because this change targets CLI tooling.
