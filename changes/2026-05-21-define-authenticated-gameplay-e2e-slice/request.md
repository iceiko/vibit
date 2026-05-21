# Request

Date: 2026-05-21
Requester: Maintainer

## Maintainer Prompt

```text
继续推进
```

## Interpreted Request

Advance the next ready work item from `.arch/work-items.yaml`:

```text
M-112/W-0184 Define authenticated gameplay E2E slice
```

Per `docs/workflow.md`, this means advance one next-ready work item unless blocked.

## Scope

Define and prove the first authenticated gameplay end-to-end path:

```text
local onboarding -> login -> connection binding -> protected inventory -> presence query -> logout
```

Use existing runtime capabilities first. Do not add protocol routes, Protobuf sources, generated output, migrations, dependencies, release artifacts, production signup, external identity providers, password login, account recovery, multi-device linking, broad product modules, or direct Nakama/Pitaya API compatibility unless the change spec explicitly scopes them.
