# Request

Implement `W-0290 Implement Pitaya-aligned startup and shutdown hook source-first map`.

The change must expose `node tools/vibit inspect pitaya-startup-shutdown --json` and register `runtime.pitaya_aligned_startup_shutdown_hook_source_first_map`.

It must report current explicit bootstrap composition and lifecycle hook deferrals without adding startup hook behavior, shutdown hook behavior, lifecycle hook execution, dependency container behavior, component discovery behavior, component loading behavior, component module loading behavior, handler module registration behavior, handler registration behavior, dynamic handler registration, runtime endpoint behavior, dashboards, admin console behavior, protocol messages or routes, Protobuf sources, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, authentication/session behavior changes, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.
