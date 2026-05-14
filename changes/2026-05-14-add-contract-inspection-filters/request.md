# Request

## Original Request

Continue advancing the generator and contract tooling hardening milestone.

## Clarified Requirement

Extend contract index inspection so agents can filter registered contracts by module, contract type, and source status before reading broad contract trees.

## User-Visible Outcome

Agents can run:

```bash
node tools/vibit inspect contracts --type command --json
node tools/vibit inspect contracts --status draft --json
node tools/vibit inspect contracts --module inventory --type query --json
```

## Non-Goals

- Do not change contract semantics.
- Do not change contract source format.
- Do not replace the contract registry.
- Do not implement runtime behavior, authentication, persistence, protocol envelope changes, WebSocket handshake changes, or player handlers.

## Acceptance Criteria

- `inspect contracts` supports `--type <contract-type>`.
- `inspect contracts` supports `--status <status>`.
- Existing unfiltered and `--module` contract inspection remains compatible.
- Agent tooling docs and translation mention the narrower inspection commands.
