# Plan

1. Review the current Nakama-first roadmap, prototype-ready maturity criteria, alpha developer flow, and friends relationship route local proof.
2. Compare current candidate directions:
   - minimum operations inspection surface;
   - groups relationship lifecycle gate;
   - parties lifecycle gate;
   - chat streams gate;
   - leaderboards or tournaments gate;
   - matchmaking gate;
   - authoritative match runtime gate;
   - SDK or generated client library publication;
   - Pitaya distributed architecture review.
3. Select the next bounded capability family and record rationale.
4. Add `ADR-0151`.
5. Mark `M-171/W-0243` completed and open `M-172/W-0244`.
6. Update runtime, reference, contracts, conventions, module, alpha, roadmap, and agent guide records.
7. Add repository check coverage for the selection.
8. Run focused and repository-wide verification.

## Selected Path

Select `define_minimum_operations_inspection_surface_gate` as the next direction.

The follow-up work item is a gate, not implementation. It should define the accepted minimum operations inspection posture before adding any endpoint, route, console, metrics, observability, or state-inspection behavior.

## Verification

Required commands:

```bash
node -c tools/vibit
node tools/vibit inspect rule runtime.next_nakama_prototype_ready_capability_after_friends_route_proof
node tools/vibit inspect next --json
node tools/vibit check change select-next-nakama-prototype-ready-capability-after-friends-route-proof --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

## Rollback

Rollback removes the W-0243 selection records, ADR-0151, rule registration, and W-0244 next-ready update. No runtime, protocol, generated output, migration, dependency, data, release, hosted, SDK, operations implementation, or direct compatibility rollback is required.
