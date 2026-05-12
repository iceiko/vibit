# Verification

Verified:

- `node tools/vibit inspect memory`
- `node tools/vibit inspect memory | node -e '<JSON.parse assertion>'`
- `node tools/vibit --help | rg "inspect memory"`
- `node tools/vibit check schemas`
- `node tools/vibit check all`
- `node tools/vibit check all --json | node -e '<JSON.parse assertion>'`
- `rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" . .git/config`
- `git diff --check`

Not verified:

- Full JSON Schema validation is not implemented yet; current schema checks validate critical fields only.

Not applicable:

- Runtime server tests do not apply because this change targets CLI tooling.
