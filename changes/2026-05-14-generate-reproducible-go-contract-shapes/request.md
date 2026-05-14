# Request

## Original Request

Continue advancing the selected generator and contract tooling hardening path.

## Clarified Requirement

Add a reproducible generator for Go contract shape files from semantic contract manifests, and extend generated-output checks to verify source trace and drift.

## User-Visible Outcome

Agents can run:

```bash
node tools/vibit generate contract-shapes all
node tools/vibit inspect generated --json
node tools/vibit check generated --json
```

The repository contains generated contract shape files for registered inventory and player semantic contracts.

## Non-Goals

- Do not generate runtime handlers.
- Do not generate WebSocket routes.
- Do not implement persistence, authentication, token, credential, or session behavior.
- Do not change Protobuf envelope or gameplay protocol behavior.

## Acceptance Criteria

- Generator command exists.
- Generated files include generated marker, source trace, generator trace, and contract trace.
- `check generated` detects missing or drifted contract shapes.
- Generated files remain metadata only.
