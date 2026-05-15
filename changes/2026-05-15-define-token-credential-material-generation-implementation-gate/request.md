# Request

Define the token and credential material generation implementation gate as the next bounded work item after the environment verifier key loader implementation.

The change must stay gate-only. It must not add token generation code, credential generation code, digest helpers, verifier comparison, authentication service behavior, protocol carriers, repository changes, migrations, startup wiring, new dependencies, or production authentication behavior.
