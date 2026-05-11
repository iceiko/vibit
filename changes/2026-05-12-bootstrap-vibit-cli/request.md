# Request

## Original Request

The maintainer approved continuing with the recommended next step: create the first executable vibit CLI prototype.

## Clarified Requirement

Create the first minimal CLI prototype that can turn the current standards into executable checks and generation entry points.

Initial target commands:

```text
vibit check architecture
vibit check module <module>
vibit check change <change-id>
vibit generate module <module>
```

The first implementation should be small and standards-oriented. It should validate current manifests and file layout before attempting a full server framework.

## User-Visible Outcome

A developer or agent can run a local CLI command to inspect the repository and receive deterministic, actionable output.

## Non-Goals

- Do not implement a full game server runtime.
- Do not generate production application code yet.
- Do not choose a final long-term implementation language without recording the decision.
- Do not add heavyweight dependencies unless justified.

## Unknowns

- Which implementation language should be used for the first CLI?
- Should the CLI be a standalone script first or a packaged executable?
- How strict should the first architecture checks be?

## Acceptance Criteria

- [ ] Implementation language selected and recorded.
- [ ] CLI can print help.
- [ ] `vibit check architecture` checks required root docs and manifests.
- [ ] `vibit check change <change-id>` checks required change spec files.
- [ ] `vibit check module <module>` reports missing module manifests clearly.
- [ ] `vibit generate module <module>` creates a minimal module skeleton or explicitly remains a documented stub.
- [ ] Verification commands are documented.
