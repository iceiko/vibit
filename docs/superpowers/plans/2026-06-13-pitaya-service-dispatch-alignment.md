# Pitaya Service Dispatch Alignment Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add the ten confirmed Pitaya-aligned service/export/dispatch/source-first planning steps without implementing distributed runtime behavior.

**Architecture:** The work extends the existing Pitaya source-first sequence after `W-0302` by adding five gate/map pairs from `W-0303` through `W-0312`. Runtime behavior remains unchanged; `tools/vibit` gains table-driven inspection and rule checks for the new source-first maps, while docs, ADRs, manifests, change specs, and continuation text record the new boundaries.

**Tech Stack:** Markdown, YAML manifests, Node.js `tools/vibit`, repository check rules, existing vibit change-spec and ADR conventions.

---

### Task 1: Source-First Metadata And Plan Artifacts

**Files:**
- Create: `docs/superpowers/plans/2026-06-13-pitaya-service-dispatch-alignment.md`
- Create: `docs/pitaya-aligned-service-export-boundary-gate.md`
- Create: `docs/pitaya-aligned-service-export-boundary-gate.zh-CN.md`
- Create: `docs/pitaya-aligned-remote-call-dispatch-boundary-gate.md`
- Create: `docs/pitaya-aligned-remote-call-dispatch-boundary-gate.zh-CN.md`
- Create: `docs/pitaya-aligned-frontend-message-forwarding-boundary-gate.md`
- Create: `docs/pitaya-aligned-frontend-message-forwarding-boundary-gate.zh-CN.md`
- Create: `docs/pitaya-aligned-backend-service-route-ownership-boundary-gate.md`
- Create: `docs/pitaya-aligned-backend-service-route-ownership-boundary-gate.zh-CN.md`
- Create: `docs/pitaya-aligned-cluster-event-bus-boundary-gate.md`
- Create: `docs/pitaya-aligned-cluster-event-bus-boundary-gate.zh-CN.md`
- Create: `decisions/ADR-0211-pitaya-aligned-service-export-boundary-gate.md` through `decisions/ADR-0220-pitaya-aligned-cluster-event-bus-source-first-map.md`
- Create: ten `changes/2026-06-13-.../` directories
- Create: ten `conversations/2026-06-13-...md` files

- [ ] **Step 1: Verify the new first command and rules fail**

Run:

```bash
node tools/vibit inspect pitaya-service-export --json
node tools/vibit inspect rule runtime.pitaya_aligned_service_export_boundary_gate
node tools/vibit inspect rule runtime.pitaya_aligned_service_export_source_first_map
```

Expected: the inspect command exits with `Unknown command`; both rule inspections exit with `Unknown rule_id`.

- [ ] **Step 2: Add the ten documentation, ADR, change, and conversation artifacts**

Use the existing Pitaya startup/shutdown and component-loading artifacts as templates, replacing the capability vocabulary, ADR numbers, work item numbers, and deferrals exactly for the ten confirmed steps.

- [ ] **Step 3: Review artifacts for forbidden scope**

Run:

```bash
rg -n "pitaya\\.|/v2/|nk\\.|dependency_added: true|distributed_runtime_implementation_added: true|protocol_route_added: true|protobuf_source_added: true|generated_output_added: true" docs/pitaya-aligned-*-boundary-gate*.md decisions/ADR-021*.md decisions/ADR-0220-*.md changes/2026-06-13-*/ conversations/2026-06-13-*.md
```

Expected: no direct compatibility markers and no true flags for the forbidden implementation surfaces.

### Task 2: Work Queue And Manifest Continuation

**Files:**
- Modify: `.arch/work-items.yaml`
- Modify: `.arch/runtime.yaml`
- Modify: `.arch/reference.yaml`
- Modify: `.arch/conventions.yaml`
- Modify: `.arch/contracts.yaml`
- Modify: `.arch/modules.yaml`
- Modify: `modules/currency/module.yaml`
- Modify: `README.md`
- Modify: `README.zh-CN.md`
- Modify: `AGENTS.md`
- Modify: `AGENTS.zh-CN.md`
- Modify: `runtime/AGENTS.md`
- Modify: `runtime/AGENTS.zh-CN.md`
- Modify: `docs/v0.1-alpha-goal.md`
- Modify: `docs/v0.1-alpha-goal.zh-CN.md`
- Modify: `docs/alpha-developer-flow.md`
- Modify: `docs/alpha-developer-flow.zh-CN.md`
- Modify: `docs/alpha-acceptance-checklist.md`
- Modify: `docs/alpha-acceptance-checklist.zh-CN.md`
- Modify: `docs/product-maturity-milestones.md`
- Modify: `docs/product-maturity-milestones.zh-CN.md`
- Modify: `docs/nakama-pitaya-product-parity-roadmap.md`
- Modify: `docs/nakama-pitaya-product-parity-roadmap.zh-CN.md`

- [ ] **Step 1: Append milestones `M-231` through `M-240` and work items `W-0303` through `W-0312`**

Expected: each work item is completed, depends on the previous item, records its ADR/check rule/change path, and the final `W-0312` leaves no next-ready item.

- [ ] **Step 2: Update repository continuation text**

Expected: README, AGENTS, runtime AGENTS, module manifests, alpha docs, and product roadmap point to the new final state: `W-0312` completed and no repository next work item currently ready.

- [ ] **Step 3: Run work check**

Run:

```bash
node tools/vibit check work --json
```

Expected: passes with `current_milestone: M-240` and no `next_ready` entries.

### Task 3: Inspection Commands And Rule Checks

**Files:**
- Modify: `tools/vibit`
- Modify: `rules/check-rules.json`

- [ ] **Step 1: Add known rule ids and rule catalog records**

Expected: the ten new rule ids are present in `KNOWN_CHECK_RULE_IDS` and `rules/check-rules.json`.

- [ ] **Step 2: Add table-driven Pitaya sequence metadata**

Expected: the table captures command names, capability keys, gate and map work item ids, ADR ids, change paths, allowed vocabulary, source surfaces, current mappings, deferrals, redaction fields, and next-step metadata.

- [ ] **Step 3: Add inspect commands**

Run:

```bash
node tools/vibit inspect pitaya-service-export --json
node tools/vibit inspect pitaya-remote-call-dispatch --json
node tools/vibit inspect pitaya-frontend-forwarding --json
node tools/vibit inspect pitaya-backend-service-routes --json
node tools/vibit inspect pitaya-cluster-event-bus --json
```

Expected: all five commands exit 0 and return source-first JSON with false implementation deferrals.

- [ ] **Step 4: Add runtime check dispatch**

Run:

```bash
node tools/vibit inspect rule runtime.pitaya_aligned_service_export_boundary_gate
node tools/vibit inspect rule runtime.pitaya_aligned_service_export_source_first_map
node tools/vibit check runtime --json
```

Expected: the rule inspections find catalog entries and runtime check passes.

### Task 4: Full Verification, Commit, And Push

**Files:**
- All files touched by Tasks 1 through 3

- [ ] **Step 1: Run syntax and targeted inspections**

Run:

```bash
node -c tools/vibit
node tools/vibit inspect pitaya-service-export --json
node tools/vibit inspect pitaya-remote-call-dispatch --json
node tools/vibit inspect pitaya-frontend-forwarding --json
node tools/vibit inspect pitaya-backend-service-routes --json
node tools/vibit inspect pitaya-cluster-event-bus --json
```

Expected: all commands exit 0.

- [ ] **Step 2: Run repository checks**

Run:

```bash
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check all >/dev/null
git diff --check
```

Expected: all commands exit 0.

- [ ] **Step 3: Run Go tests**

Run:

```bash
cd runtime && go test ./...
```

Expected: all packages pass.

- [ ] **Step 4: Commit and push**

Run:

```bash
git status --short --ignored
git add .
git commit -m "Add Pitaya service dispatch source-first maps"
```

Expected: ignored `.vibit.local.env` and `node_modules/` are not staged. Push using only the ignored `GITHUB_TOKEN` variable from `.vibit.local.env`, without printing the token value.
