# Request

## Original Request

The maintainer asked to continue after the release and previously stated that the README should help vibit find users.

## Clarified Requirement

Define `W-0196 Define first alpha user discovery loop` as a bounded documentation, ADR, manifest, and check-rule slice. The work must record target developer segments, outreach surfaces, feedback capture, success signals, and stop conditions after the source-first `v0.1.0-alpha.1` release.

## User-Visible Outcome

Maintainers and agents can now inspect `docs/first-alpha-user-discovery-loop.md` to understand how vibit should find first alpha users without executing a broad announcement or expanding release/runtime scope.

## Non-Goals

- Execute public announcements beyond the existing GitHub release record.
- Run paid promotion.
- Add hosted deployments or demos.
- Create binaries, packages, containers, checksums, provenance files, signing artifacts, install scripts, registries, or SDK packages.
- Change runtime behavior.
- Add protocol routes, Protobuf sources, generated output, migrations, dependencies, broad operations/admin behavior, authentication/session behavior, broad product modules, or direct Nakama/Pitaya API compatibility.

## Unknowns

- Which feedback intake surface should be prepared first: GitHub issues, GitHub discussions, labels, templates, or a smaller repository-owned note.
- Whether the maintainer wants broad public outreach after the intake surface exists.

## Acceptance Criteria

- [x] Target developer segments are recorded.
- [x] Outreach surfaces are recorded without executing announcements.
- [x] Feedback capture fields and review questions are recorded.
- [x] Success signals and stop conditions are recorded.
- [x] Next bounded work item is set to feedback intake surface preparation.
- [x] Runtime, protocol, generated-output, migration, dependency, operations/admin, authentication/session, hosted deployment, broad artifact, product module, and direct compatibility deferrals remain preserved.
