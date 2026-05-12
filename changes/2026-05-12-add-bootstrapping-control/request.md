# Request

## Original Request

The maintainer shifted the focus from whether the project has already drifted to preventing future drift:

> 比起现在是否已经跑偏，我更关注的是未来是否可能跑偏。我们要让整个项目处在一个自举可控的范围内，这样才能最终很好地达成我们的目标。

The maintainer also provided a GitHub token and asked the agent to commit and push the current work first.

## Clarified Requirement

After pushing the completed local commits, add a small governance rule that keeps vibit self-bootstrapping productive and prevents future meta-tooling drift.

## User-Visible Outcome

Future agents have a clear rule: new standards, inspect commands, check commands, schemas, generators, and workflow rules should directly support concrete runtime progress, module boundaries, contracts, tests, verification, or agent context reduction.

## Non-Goals

- Do not add a new CLI command for this rule yet.
- Do not block all future tooling work.
- Do not decide the runtime language in this change.
- Do not start the first runtime vertical slice in this change.

## Unknowns

- Whether this rule should later become an automated check.
- Which runtime language and architecture will be chosen first.
- Which backend vertical slice will be implemented first.

## Acceptance Criteria

- [x] Previously completed local commits are pushed.
- [x] The constitution includes a bootstrapping control principle.
- [x] AGENTS includes an operational bootstrapping control rule.
- [x] An ADR records the durable decision.
- [x] Simplified Chinese translations are updated.
- [x] Conversation log records the maintainer's concern.
- [x] Verification is recorded.
