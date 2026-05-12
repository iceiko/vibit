# Verification

Verified:

- `node tools/vibit inspect rule check.subcheck`
- `node tools/vibit inspect rule check.subcheck | node -e '<JSON.parse assertion>'`
- `node tools/vibit inspect rule missing.rule`
- `node tools/vibit check all`
- `node tools/vibit check all --json`
- `rg -n "ghp_[A-Za-z0-9]|github_pat_[A-Za-z0-9]" . .git/config`

Not verified:

- Full JSON Schema validation is not implemented yet; current schema checks validate critical fields only.

Not applicable:

- Runtime server tests do not apply because this change targets CLI tooling.
