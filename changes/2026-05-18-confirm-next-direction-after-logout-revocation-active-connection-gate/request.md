# Request

Select the next bounded milestone direction after the logout/revocation active-connection gate.

The maintainer asked the agent in Chinese to recommend the next ten steps and continue, while prioritizing lessons from Nakama and Pitaya.

Selected direction:

```text
define_logout_access_token_behavior_gate
```

Reason:

- `W-0154` recommended `presented_access_token_only` as the first future logout scope.
- Before code revokes tokens, vibit needs a precise behavior gate for proof validation, verifier comparison, transaction ordering, public error collapse, and explicit session/socket/protocol deferrals.
- Nakama reinforces authentication/session lifecycle pressure.
- Pitaya reinforces keeping session/connection infrastructure separate from handlers and credential parsing.
