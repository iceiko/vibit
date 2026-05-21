# Request

## Original Request

The maintainer asked to turn the agent's maturity assessment into project milestones, preserve that direction in the repository, and continue pushing vibit toward a real product and productivity stage.

## Clarified Requirement

Complete `W-0197 Prepare first alpha feedback intake surfaces` by preparing a repository-owned early alpha feedback intake surface, recording redaction-safe triage guidance, and preserving the product maturity path from source-first alpha to prototype-ready foundation, single-node production-candidate foundation, and Nakama/Pitaya-class product.

## User-Visible Outcome

Maintainers, contributors, and early users can now:

- open structured alpha feedback through `.github/ISSUE_TEMPLATE/alpha-feedback.yml`;
- understand how feedback is mapped to product maturity stages;
- see that `v0.1.0-alpha.1` is a source-first alpha, not production-ready;
- see that the next product stage is prototype-ready foundation;
- continue with `W-0198 Define prototype-ready foundation execution plan`.

## Non-Goals

- Do not execute public announcements beyond the GitHub release record.
- Do not run paid promotion.
- Do not add hosted deployments or demos.
- Do not create release binaries, packages, containers, checksums, provenance files, signing artifacts, install scripts, registry publications, or SDK packages.
- Do not change runtime behavior.
- Do not add protocol routes or generated output.
- Do not add migrations or dependencies.
- Do not broaden operations/admin behavior.
- Do not change authentication/session behavior.
- Do not add broad product modules.
- Do not add direct Nakama/Pitaya API compatibility.

## Unknowns

- Whether GitHub Discussions should be added later remains undecided.
- Which prototype-ready capability should lead after the execution plan remains undecided.
- Whether broader public outreach should happen remains a later maintainer authorization decision.

## Acceptance Criteria

- [x] A repository-owned feedback intake surface exists.
- [x] Feedback fields, redaction expectations, triage labels, and triage posture are recorded.
- [x] Product maturity stages are recorded in durable documentation.
- [x] `.arch` manifests and README pointers advance to the next work item.
- [x] Repository checks cover the new feedback-intake state.
- [x] Runtime, protocol, generated output, migration, dependency, release artifact, hosted deployment, public announcement, paid promotion, authentication/session, broad product, and direct compatibility deferrals remain preserved.
