# Request

Implement `W-0284 Implement Pitaya-aligned handler module registration source-first map`.

The slice must expose a source-first inspection command for handler module registration vocabulary and current repository mappings after `ADR-0191`, without implementing handler module registration behavior, handler registration behavior, dynamic handler registration, component discovery or loading, component module loading, startup hooks, shutdown hooks, runtime endpoint behavior, dashboards, admin console behavior, protocol shape, generated output, persistence, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.

The requested inspection command is:

```sh
node tools/vibit inspect pitaya-handler-modules --json
```

The repository check rule is:

```text
runtime.pitaya_aligned_handler_module_registration_source_first_map
```

The follow-up work item opened by this slice is `W-0285 Select next Pitaya-aligned direction after handler module registration map`.
