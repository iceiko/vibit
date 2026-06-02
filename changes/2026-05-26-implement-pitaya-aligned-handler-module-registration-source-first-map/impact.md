# Impact

## Included

- Add `runtime.pitaya_aligned_handler_module_registration_source_first_map` to the rule catalog and runtime checks.
- Add `node tools/vibit inspect pitaya-handler-modules --json`.
- Report handler module registration vocabulary, source-first mapping, source surfaces, redaction posture, and implementation deferrals.
- Record `ADR-0192`, the change artifacts, and a conversation log.
- Mark `M-212/W-0284` completed and open `M-213/W-0285` as next-ready.
- Update repository continuation memory in manifests, README files, AGENTS guides, and roadmap documents.

## Excluded

- Handler module registration behavior.
- Handler registration behavior.
- Dynamic handler registration.
- Component discovery or loading.
- Component module loading.
- Startup hooks or shutdown hooks.
- Runtime endpoint behavior.
- Dashboard or admin console behavior.
- Protocol messages, routes, Protobuf source, or generated output.
- Repository interfaces, PostgreSQL adapters, migrations, persistence, or dependencies.
- Authentication/session behavior changes.
- Hosted deployment, SDK publication, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility.
