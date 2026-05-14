# Request

## Original Request

它区域限制呗。按照你的建议和判断推进10步，除非有非常必要的，需要我决策的，再停下来问。

## Clarified Requirement

Advance `W-0079` by hardening local repository checks for the ratified credential and token verifier migration sources.

## User-Visible Outcome

Maintainers and agents can detect authentication migration drift earlier through local checks, before repository interfaces or adapters exist.

## Non-Goals

- Do not add authentication repository interfaces.
- Do not add PostgreSQL authentication adapters.
- Do not implement runtime credential lookup, token issuance, token validation, logout, refresh, cleanup, handlers, routes, or production authentication behavior.
- Do not weaken existing runtime, protocol, or selected login/token checks.
- Do not require live PostgreSQL for the default check path.

## Unknowns

- The exact shape of the static check additions may evolve after reading the current tool rules.
- Whether additional documentation references need to be updated will depend on the final tool logic.

## Acceptance Criteria

- [ ] Add local static checks for the credential and token verifier migration sources.
- [ ] Keep repository-relative JSON output using forward slashes on every platform.
- [ ] Preserve player account lifecycle table separation.
- [ ] Preserve the deferral of repositories, adapters, and runtime authentication behavior.
