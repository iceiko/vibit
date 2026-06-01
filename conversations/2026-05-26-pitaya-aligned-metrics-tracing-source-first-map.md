# Conversation: Pitaya-Aligned Metrics And Tracing Source-First Map

Date: 2026-06-01
Work item: W-0275
Decision: ADR-0183
Rule: runtime.pitaya_aligned_metrics_tracing_source_first_map

## Context

The maintainer asked to continue advancing toward Pitaya alignment in bounded steps and to commit and push completed increments.

W-0274 completed the gate-only metrics and tracing boundary, accepted ADR-0182, and opened W-0275 as the source-first inspection map follow-up.

## Maintainer Narrative

Continue toward Pitaya alignment while preserving source-first, checkable boundaries and avoiding runtime behavior unless a bounded work item authorizes it.

## Agent Response Summary

The agent implemented W-0275 as a source-first tooling map. The change adds `node tools/vibit inspect pitaya-metrics-tracing --json`, registers a runtime check rule, records ADR-0183, and opens W-0276 as the next direction selection work item.

## RED Checks

```text
node tools/vibit inspect rule runtime.pitaya_aligned_metrics_tracing_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_metrics_tracing_source_first_map

node tools/vibit inspect pitaya-metrics-tracing --json
# Unknown command.

node tools/vibit check change implement-pitaya-aligned-metrics-tracing-source-first-map --json
# change directory does not exist
```

## Decisions

Implement a source-first Pitaya-aligned metrics and tracing inspection map that reports vocabulary, current runtime observability mapping, operations posture, trace/correlation/sampling deferrals, redaction posture, source surfaces, and deferrals without adding runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, dependencies, hosted surfaces, or direct compatibility.

## Artifacts

- `decisions/ADR-0183-pitaya-aligned-metrics-tracing-source-first-map.md`
- `changes/2026-05-26-implement-pitaya-aligned-metrics-tracing-source-first-map/`
- `tools/vibit`
- `rules/check-rules.json`
- `.arch/work-items.yaml`
- `.arch/runtime.yaml`
- `.arch/reference.yaml`
- `.arch/conventions.yaml`
- `.arch/contracts.yaml`

## Boundaries Preserved

No runtime endpoint behavior, metrics endpoints, tracing pipelines, observability pipelines, dashboards, admin console behavior, player/session/token inspectors, event/audit tables, protocol messages, Protobuf sources, generated output, repository interfaces, PostgreSQL adapters, migrations, dependencies, hosted surfaces, SDKs, release artifacts, distributed runtime behavior, or direct Nakama/Pitaya API compatibility were added.

No ignored credential file contents were read or printed.

## Follow-Up

Open W-0276 as the next Pitaya-aligned direction selection after the metrics and tracing map.

## Open Questions

- Whether the next Pitaya direction should continue architecture vocabulary coverage or pivot back to Nakama-facing product capability families.
- Whether later metrics, tracing, dashboard, admin console, correlation, sampling, and distributed node telemetry gates should remain separate.

## Redaction Notes

The command output is source-first metadata only. It does not expose raw credentials, raw tokens, lookup digests, verifier digests, verifier keys, DSNs, headers, cookies, query strings, WebSocket subprotocol values, remote addresses, database payloads, local secret file contents, route payloads, session data payloads, concrete transport metadata, or ignored credential file contents.
