# Request

Define the process environment verifier key loader gate without implementing environment parsing, Base64 text decoding, startup wiring, token generation, credential generation, verifier digest computation, verifier comparison, authentication service behavior, protocol carriers, repository changes, migration changes, external secret-manager integrations, dependencies, or production authentication behavior.

This change advances `W-0098` under `M-026`.

## Maintainer Intent

The maintainer asked the agent to continue bounded work without unnecessary confirmation. `W-0098` is the next ready work item, and `W-0097` has already implemented the explicit in-memory verifier key set validator that the future loader must call.

## Required Outcome

- Add the English and Simplified Chinese environment verifier key loader gate standard.
- Add ADR-0046.
- Declare the future environment variable names.
- Declare the Base64 decoding posture.
- Require future loader handoff to `NewVerifierKeySet`.
- Preserve deferrals for loader implementation, startup wiring, secret files, `.env`, KMS, cloud secret managers, token behavior, digest helpers, service behavior, protocol carriers, repositories, migrations, dependencies, and production behavior.
