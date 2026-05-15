# Request

Define the verifier digest helper implementation gate as the next bounded work item after the token and credential material generation helper implementation.

The change must stay gate-only. It must not add digest computation code, HMAC helpers, verifier comparison, authentication service behavior, protocol carriers, repository changes, migrations, startup wiring, new dependencies, or production authentication behavior.
