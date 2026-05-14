# Request

## Original Request

继续推进十步，除非需要遇到我确认的问题，否则不要停下来。

## Clarified Requirement

Advance `W-0087` by closing `M-016 Runtime Authentication Implementation Boundary Planning` after the runtime authentication implementation boundary is defined, and open the next bounded preparation milestone without implementing runtime authentication behavior.

## User-Visible Outcome

Maintainers and future agents should see:

- `M-016` as completed.
- `W-0087` as completed.
- `M-017 Authentication Generated Contract Shape Gate` as active.
- `W-0088 Decide authentication generated contract shape timing` as the next ready work item.

## Non-Goals

- Do not implement runtime authentication behavior.
- Do not generate authentication contract shapes in this closeout.
- Do not add token generation, verifier comparison, token validation, login execution, logout execution, refresh, or cleanup jobs.
- Do not add Protobuf authentication messages or WebSocket proof carriers.
- Do not add authentication dependencies.
- Do not change `authentication.Repository`.
- Do not change ratified migration schemas.

## Unknowns

- Exact generated authentication shape timing remains the next work item.
- Exact generator and check changes remain future work.

## Acceptance Criteria

- [x] Mark `M-016` completed with a completion summary.
- [x] Mark `W-0087` completed with a completion summary.
- [x] Open `M-017 Authentication Generated Contract Shape Gate`.
- [x] Make `W-0088` the only next ready work item.
- [x] Preserve runtime authentication behavior, generated output, Protobuf, WebSocket, dependency, repository, and migration deferrals.
