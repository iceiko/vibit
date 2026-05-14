# Request

## Original Request

Continue advancing the project under the selected generator and contract tooling hardening direction.

## Clarified Requirement

Add an agent-facing tooling standard and narrow JSON inspection commands for the next work item, registered contracts, generated output, and Nakama/Pitaya reference planning context.

## User-Visible Outcome

Agents can run focused inspection commands before broad source reading or edits:

- `node tools/vibit inspect next --json`
- `node tools/vibit inspect contracts --json`
- `node tools/vibit inspect generated --json`
- `node tools/vibit inspect reference --json`

## Non-Goals

- Do not add server runtime behavior.
- Do not change public gameplay protocol behavior.
- Do not implement authentication, token, credential, persistence, or player handlers.
- Do not add a major CLI dependency.

## Acceptance Criteria

- `docs/agent-tooling.md` and `docs/agent-tooling.zh-CN.md` exist.
- The new inspect commands return JSON.
- `node tools/vibit check agent-tooling --json` passes.
- `check all` includes `check agent-tooling`.
