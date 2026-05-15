# Request

Implement the explicit in-memory verifier key set validator under `runtime/internal/app/authentication` without adding environment parsing, Base64 text decoding, startup wiring, token generation, credential generation, verifier digest computation, verifier comparison, authentication service behavior, protocol carriers, repository changes, migration changes, external secret-manager integrations, dependencies, or production authentication behavior.

This change advances `W-0097` under `M-025`.

## Maintainer Intent

The maintainer asked the agent to continue bounded work without unnecessary confirmation. `W-0097` is the next ready work item and ADR-0045 defines the implementation boundary.

## Required Outcome

- Add a small explicit in-memory verifier key set validator.
- Copy input key bytes before storage.
- Avoid exposing mutable internal slices.
- Reject missing key set id, missing keys, short keys, duplicate logical keys, all-zero keys, and repeated single-byte keys.
- Keep errors redacted.
- Add focused unit tests.
