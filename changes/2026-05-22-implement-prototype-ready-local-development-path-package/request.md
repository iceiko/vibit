# Request

## Original Request

The maintainer asked to continue advancing vibit after the prototype-ready local development path gate was recorded, with the broader product direction of pushing toward a real product and production-useful foundation.

## Clarified Requirement

Complete `W-0200 Implement prototype-ready local development path package` by packaging the local development path inside the `ADR-0107` gate. The package should improve setup, startup, migration, configuration, secret redaction, example-flow, and verification ergonomics using docs, examples, static checks, and focused verification of existing behavior.

## User-Visible Outcome

Maintainers, contributors, and agents can now see:

- the prototype-ready local development path package standard;
- the fastest source checkout verification path;
- supported local prerequisites;
- a redacted placeholder local environment template;
- `.gitignore` protection for private local env files;
- explicit migration and runtime startup expectations;
- the authenticated request-loop proof and what it covers;
- local status surfaces and redaction posture;
- required verification commands;
- that `W-0201 Define storage objects behavior gate` is the next work item.

## Non-Goals

- Do not implement production runtime behavior.
- Do not add protocol routes or generated output.
- Do not add Protobuf source files.
- Do not add migrations or dependencies.
- Do not add automatic startup migration behavior.
- Do not broaden operations/admin behavior.
- Do not change authentication/session behavior.
- Do not add broad product modules.
- Do not add direct Nakama/Pitaya API compatibility.
- Do not add hosted deployments or demos.
- Do not create release binaries, packages, containers, checksums, provenance files, signing artifacts, install scripts, registry publications, SDK packages, or additional release artifacts.
- Do not execute public announcements beyond the GitHub release record.
- Do not run paid promotion.
- Do not disclose, inspect, or commit local secrets.

## Unknowns

- The first storage-object behavior is not selected by this change; it is deferred to `W-0201`.
- A future live local process client or CLI for local onboarding remains deferred.
- Containerized setup and hosted demos remain deferred until explicitly authorized.

## Acceptance Criteria

- [x] The prototype-ready local development path package is recorded in English and Simplified Chinese.
- [x] The package records the fast source path, prerequisites, redacted configuration, migration expectations, startup guidance, example flow, verification, and stop conditions.
- [x] A placeholder local env template exists without real secrets.
- [x] Private local env filenames are ignored.
- [x] Contributor entry points link the package.
- [x] Repository checks cover the new package state.
- [x] The next work item advances to `W-0201 Define storage objects behavior gate`.
- [x] Runtime, protocol, generated output, migration, dependency, release artifact, hosted deployment, public announcement, paid promotion, authentication/session, broad product, and direct compatibility deferrals remain preserved.
