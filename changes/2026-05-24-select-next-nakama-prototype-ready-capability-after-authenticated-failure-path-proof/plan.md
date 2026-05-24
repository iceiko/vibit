# Plan

1. Review the current Nakama-first roadmap, prototype-ready maturity criteria, alpha developer flow, and authenticated failure-path proof.
2. Compare current candidate directions:
   - clearer example client or example app path;
   - minimal operations inspection surface;
   - realtime delivery proof;
   - chat/group messaging gate;
   - social, competitive, matchmaking, match runtime, or distributed runtime gates.
3. Select the next bounded capability family and record rationale.
4. Add `ADR-0132`.
5. Mark `M-152/W-0224` completed and open `M-153/W-0225`.
6. Update runtime, reference, contracts, conventions, module, alpha, roadmap, and agent guide records.
7. Add repository check coverage for the selection.
8. Run focused and repository-wide verification.

## Selected Path

Select `define_local_alpha_example_client_path_gate` as the next direction.

The follow-up work item is a gate, not implementation. It should define the accepted example-client/app shape before adding files that look like a client surface.

## Verification

Required commands:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.next_nakama_prototype_ready_capability_selection
node tools/vibit check change select-next-nakama-prototype-ready-capability-after-authenticated-failure-path-proof --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## Rollback

Rollback removes the W-0224 selection records, ADR-0132, rule registration, and W-0225 next-ready update. No runtime, protocol, generated output, migration, dependency, data, release, hosted, SDK, or direct compatibility rollback is required.
