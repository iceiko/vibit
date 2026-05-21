# Plan

1. Record product maturity milestones in English and Simplified Chinese.
2. Record first alpha feedback intake guidance in English and Simplified Chinese.
3. Add a GitHub issue form for structured alpha feedback.
4. Record `ADR-0105` and conversation memory.
5. Update `.arch` manifests to complete `W-0197` and open `W-0198`.
6. Update README, AGENTS guides, alpha goal, alpha developer flow, and acceptance checklist pointers.
7. Add repository check coverage for `runtime.first_alpha_feedback_intake_surfaces`.
8. Run required verification.

## Boundary

This is a docs, workflow, feedback-intake, and check-rule slice. It must not modify runtime behavior, protocol routes, Protobuf sources, generated output, migrations, dependencies, release artifacts, hosted deployment, authentication/session behavior, broad operations/admin behavior, broad product modules, or direct Nakama/Pitaya compatibility.
