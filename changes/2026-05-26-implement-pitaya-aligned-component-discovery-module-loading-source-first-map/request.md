# Request

Implement the Pitaya-aligned component discovery and module loading source-first map for `W-0287`.

The change must expose `node tools/vibit inspect pitaya-component-loading --json` and register `runtime.pitaya_aligned_component_discovery_module_loading_source_first_map`.

The slice must not add component discovery behavior, component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.
