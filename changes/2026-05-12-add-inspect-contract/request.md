# Request

## Original Request

The maintainer asked the agent to continue the next step of the project:

> 继续下一步工作。

## Clarified Requirement

After adding `check contracts`, add a small agent-friendly inspection command for reading one registered contract as structured JSON.

Target command:

```bash
node tools/vibit inspect contract --module inventory --type command --id GrantItem
```

The command should help agents during intake by reporting:

- The matching entry in `.arch/contracts.yaml`
- The referenced contract source file
- The key source fields
- Whether the module manifest declares and references the contract
- Consistency status across registry, source, and module manifest

## User-Visible Outcome

Agents can inspect a single command, query, event, error catalog, or permission catalog without manually opening `.arch/contracts.yaml`, `contracts/<module>/...`, and `modules/<module>/module.yaml`.

## Non-Goals

- Do not add a YAML parser dependency.
- Do not implement full payload schema validation.
- Do not generate runtime code.
- Do not change existing contract source format.
- Do not replace `check contracts`.

## Unknowns

- Whether future versions should support listing all contracts for a module.
- Whether full payload schema extraction should wait for a formal YAML parser.

## Acceptance Criteria

- [x] `node tools/vibit inspect contract --module inventory --type command --id GrantItem` returns valid JSON.
- [x] The output includes registry, source, module manifest, and consistency information.
- [x] The command supports command, query, event, error, and permission contract types.
- [x] Missing contract IDs fail clearly.
- [x] The inspect output schema includes `contract_inspection`.
- [x] README and AGENTS mention the command in English and Simplified Chinese.
- [x] Verification is recorded.
