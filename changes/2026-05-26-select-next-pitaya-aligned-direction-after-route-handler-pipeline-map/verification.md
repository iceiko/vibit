# Verification

Required commands:

```bash
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect rule runtime.next_pitaya_aligned_direction_after_route_handler_pipeline_map
node tools/vibit check change select-next-pitaya-aligned-direction-after-route-handler-pipeline-map --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Status: Verified.
