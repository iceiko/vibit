# Verification

Status: Verified

Commands:

```sh
node -c tools/vibit
node tools/vibit inspect rule runtime.pitaya_aligned_backend_service_route_ownership_source_first_map
node tools/vibit inspect pitaya-backend-service-routes --json
node tools/vibit check change implement-pitaya-aligned-backend-service-route-ownership-source-first-map --json
node tools/vibit check runtime --json
node tools/vibit check work --json
```

Focused TDD evidence:

- Before this change, `node tools/vibit inspect rule runtime.pitaya_aligned_backend_service_route_ownership_source_first_map` failed with `Unknown rule_id: runtime.pitaya_aligned_backend_service_route_ownership_source_first_map`.
- Before this change, `node tools/vibit inspect pitaya-backend-service-routes --json` failed with `Unknown command`.
- After implementation, the targeted command is expected to pass as part of repository verification.
