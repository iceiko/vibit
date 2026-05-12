# Verification

Verified:

- `node tools/vibit inspect contract --module inventory --type command --id GrantItem`
  - Result: passed; returned a consistent `contract_inspection` for `GrantItem`.
- `node tools/vibit generate contract --module inventory --type command --id GrantItem`
  - Result: passed; generated `modules/inventory/generated/contracts/GrantItem.generated.ts` from `contracts/inventory/commands/GrantItem.yaml`.
- `node tools/vibit check generated`
  - Result: passed with 6 OK checks and 0 warnings.
- `node tools/vibit check generated --json`
  - Result: passed with `status: "passed"`, 6 passed checks, 0 warnings, and 0 failures.
- `node tools/vibit inspect rules --category generated`
  - Result: passed; returned 3 generated-file rule metadata entries.
- `node tools/vibit check contracts --json`
  - Result: passed with 54 passed checks, 0 warnings, and 0 failures.
- `node tools/vibit check schemas`
  - Result: passed; schema files parse as JSON and the rule catalog covers the new `generated.*` rules.
- `node tools/vibit check all --json`
  - Result: passed with 27 subchecks, 27 passed, 0 warnings, and 0 failures. The subchecks include `check generated`.
- `rg -n "@generated|Source: contracts/inventory/commands/GrantItem.yaml|Generator: tools/vibit generate contract" modules/inventory/generated/contracts/GrantItem.generated.ts`
  - Result: passed; generated, source, and generator trace markers are present.
- `rg -n -o "ghp_[A-Za-z0-9]+|github_pat_[A-Za-z0-9_]+" --hidden --glob '!.git/**' --glob '!.vibit.local.env' .`
  - Result: no committed-file GitHub token pattern was found.
- `git diff --check`
  - Result: passed with no whitespace errors.

Not verified:

- Full YAML parsing is not implemented yet; this generator intentionally extracts only the fields needed for the first command shape.
- TypeScript compilation is not implemented yet because the repository does not define a TypeScript package.

Not applicable:

- TypeScript compilation is not applicable because this change does not introduce a TypeScript package, package manager, or compiler.
- Runtime server tests are not applicable because this change does not implement runtime behavior.
