# Request

Define the local verifier key configuration loading gate before adding secret-loading code, environment parsing, token generation, credential generation, verifier digest computation, verifier comparison, authentication service behavior, protocol carriers, repository changes, migration changes, dependencies, or production authentication behavior.

This change advances `W-0096` under `M-024`.

## Maintainer Intent

The maintainer asked the agent to continue the work queue and expects bounded professional decisions to proceed without unnecessary confirmation.

The project should keep doing necessary preparation before implementation accelerates. Authentication work should remain controlled by explicit gates and repository checks.

## Required Outcome

- Define the first local verifier key configuration loading code gate.
- Choose the first implementation sequence.
- Define owner package, allowed files, forbidden files, validation rules, redaction rules, required tests, dependency posture, and deferrals.
- Preserve all runtime authentication behavior deferrals.
