# Verification

Verified:

- `node tools/vibit inspect contract --module inventory --type command --id GrantItem`
  - Result: passed; returned `contract_inspection` JSON with registry, source, module manifest, and consistency fields.
- `node tools/vibit inspect contract --module inventory --type command --id GrantItem | node -e '<JSON.parse assertion>'`
  - Result: passed; output kind, contract id, consistency status, and source path matched expectations.
- `node tools/vibit inspect contract --module inventory --type error --id inventory_errors | node -e '<JSON.parse assertion>'`
  - Result: passed; returned 3 error codes.
- `node tools/vibit inspect contract --module inventory --type permission --id inventory_permissions | node -e '<JSON.parse assertion>'`
  - Result: passed; returned 2 permission ids.
- `node tools/vibit inspect contract --module inventory --type command --id MissingCommand`
  - Result: failed as expected with a clear unknown contract message.
- `node tools/vibit --help | rg "inspect contract"`
  - Result: passed; help output includes the command and description.
- `node tools/vibit check contracts --json`
  - Result: passed with 54 passed checks, 0 warnings, and 0 failures.
- `node tools/vibit check schemas`
  - Result: passed; schema files parse as JSON and critical schema/catalog fields validate.
- `node tools/vibit check all --json`
  - Result: passed with 25 subchecks, 25 passed, 0 warnings, and 0 failures.
- `rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" --hidden --glob '!.git/**' --glob '!.vibit.local.env' .`
  - Result: no committed-file GitHub token pattern was found.
- `git diff --check`
  - Result: passed with no whitespace errors.

Not verified:

- Full YAML parsing is not implemented yet; this command intentionally uses lightweight field extraction.
- Full payload schema validation is not implemented yet.

Not applicable:

- Runtime server tests do not apply because this change targets CLI tooling.
