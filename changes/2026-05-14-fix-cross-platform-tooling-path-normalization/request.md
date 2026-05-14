# Request

## Original Request

The maintainer reported two cross-platform tooling bugs from another AI review:

```text
BUG-001: Path boundary checks falsely report legal files on Windows because path prefix checks compare forward-slash prefixes against Windows backslash paths.
BUG-002: JSON output artifact paths differ by platform because Windows emits backslashes while macOS/Linux emit forward slashes.
```

The maintainer asked to make the project compatible and add Windows platform support to the project standard.

## Clarified Requirement

Normalize repository-relative paths used by tooling to forward-slash form before JSON output and before repository path prefix comparisons.

Document the cross-platform JSON path rule in the agent tooling standard and paired Simplified Chinese translation.

## User-Visible Outcome

Agents consuming `node tools/vibit ... --json` should see stable forward-slash repository paths in `artifact` fields regardless of operating system.

Runtime boundary checks should not misclassify valid Windows paths as outside allowed owner roots.

## Non-Goals

- Do not change runtime server behavior.
- Do not add Windows-specific runtime support beyond tooling path normalization.
- Do not introduce external path libraries.
- Do not change public game protocol, contracts, migrations, authentication behavior, or generated output conventions.

## Unknowns

Windows execution is not available in the current Termux environment, so verification is by code inspection and Linux/Termux command execution.

## Acceptance Criteria

- [x] `rel()` returns forward-slash repository-relative paths.
- [x] JSON check result `artifact` values are normalized to forward-slash paths.
- [x] Runtime boundary prefix checks use a cross-platform helper instead of raw path `startsWith` on repository paths.
- [x] Agent tooling standard documents cross-platform JSON path normalization.
- [x] Simplified Chinese translation is updated.
- [x] Verification is recorded.
