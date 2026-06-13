# Verification

Status: Verified

Commands:

```sh
node -c tools/vibit
node tools/vibit inspect rule runtime.pitaya_aligned_cluster_event_bus_source_first_map
node tools/vibit inspect pitaya-cluster-event-bus --json
node tools/vibit check change implement-pitaya-aligned-cluster-event-bus-source-first-map --json
node tools/vibit check runtime --json
node tools/vibit check work --json
```

Focused TDD evidence:

- Before this change, `node tools/vibit inspect rule runtime.pitaya_aligned_cluster_event_bus_source_first_map` failed with `Unknown rule_id: runtime.pitaya_aligned_cluster_event_bus_source_first_map`.
- Before this change, `node tools/vibit inspect pitaya-cluster-event-bus --json` failed with `Unknown command`.
- After implementation, the targeted command is expected to pass as part of repository verification.
