# Pitaya Distributed Operations Alignment Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add twenty confirmed Pitaya-aligned distributed operations source-first planning steps without implementing distributed runtime behavior.

**Architecture:** Extend the completed Pitaya service dispatch sequence after `W-0312` with ten gate/map pairs from `W-0313` through `W-0332`. Runtime behavior remains unchanged; `tools/vibit` gains table-driven inspection and rule checks for the new source-first maps, while docs, ADRs, manifests, change specs, and continuation text record the new boundaries.

**Tech Stack:** Markdown, YAML manifests, Node.js `tools/vibit`, repository check rules, existing vibit change-spec and ADR conventions.

---

### Task 1: RED Baseline

**Files:**
- Modify later: `tools/vibit`
- Modify later: `rules/check-rules.json`

- [ ] **Step 1: Verify first new inspect command is absent**

Run:

```bash
node tools/vibit inspect pitaya-node-identity --json
```

Expected: exits non-zero with `Unknown command`.

- [ ] **Step 2: Verify first new rules are absent**

Run:

```bash
node tools/vibit inspect rule runtime.pitaya_aligned_node_identity_boundary_gate
node tools/vibit inspect rule runtime.pitaya_aligned_node_identity_source_first_map
```

Expected: both exit non-zero with `Unknown rule_id`.

### Task 2: Source Artifacts

**Files:**
- Create: ten `docs/pitaya-aligned-*-boundary-gate.md`
- Create: ten paired `docs/pitaya-aligned-*-boundary-gate.zh-CN.md`
- Create: `decisions/ADR-0221-*.md` through `decisions/ADR-0240-*.md`
- Create: twenty `changes/2026-06-14-*/` directories
- Create: twenty `conversations/2026-06-14-*.md` files

- [ ] **Step 1: Generate the twenty source-first artifacts**

Use the completed service dispatch artifacts as templates. Preserve explicit false deferrals for runtime behavior, protocol routes, Protobuf source, generated output, migrations, repositories, dependencies, distributed runtime implementation, and direct compatibility.

- [ ] **Step 2: Run memory/schema checks**

Run:

```bash
node tools/vibit check schemas --json
node tools/vibit check memory --json
```

Expected: both pass.

### Task 3: Work Queue And Manifests

**Files:**
- Modify: `.arch/work-items.yaml`
- Modify: `.arch/runtime.yaml`
- Modify: `.arch/reference.yaml`
- Modify: `.arch/conventions.yaml`
- Modify: `.arch/contracts.yaml`
- Modify: `.arch/modules.yaml`
- Modify: `modules/currency/module.yaml`
- Modify: `modules/friends/module.yaml`
- Modify: `modules/storage/module.yaml`
- Modify: `README.md`
- Modify: `README.zh-CN.md`
- Modify: `AGENTS.md`
- Modify: `AGENTS.zh-CN.md`
- Modify: `runtime/AGENTS.md`
- Modify: `runtime/AGENTS.zh-CN.md`
- Modify: alpha and roadmap docs

- [ ] **Step 1: Add milestones and work items**

Add `M-241` through `M-260` and completed `W-0313` through `W-0332`. The final work item leaves no repository next-ready item.

- [ ] **Step 2: Update continuation text**

Update README, AGENTS guides, module manifests, alpha docs, and roadmap docs to point to `W-0332` / `ADR-0240`.

- [ ] **Step 3: Run work check**

Run:

```bash
node tools/vibit check work --json
```

Expected: passes with `current_milestone: M-260` and no `next_ready` entries.

### Task 4: Tooling And Rules

**Files:**
- Modify: `tools/vibit`
- Modify: `rules/check-rules.json`

- [ ] **Step 1: Add rule catalog entries**

Add the twenty new rule IDs to `KNOWN_CHECK_RULE_IDS` and `rules/check-rules.json`.

- [ ] **Step 2: Add inspect commands**

Add table-driven inspection commands:

```bash
node tools/vibit inspect pitaya-node-identity --json
node tools/vibit inspect pitaya-service-registry --json
node tools/vibit inspect pitaya-service-selector --json
node tools/vibit inspect pitaya-heartbeat-liveness --json
node tools/vibit inspect pitaya-route-targeting --json
node tools/vibit inspect pitaya-remote-timeout-retry --json
node tools/vibit inspect pitaya-distributed-session-ownership --json
node tools/vibit inspect pitaya-presence-fanout --json
node tools/vibit inspect pitaya-cross-node-errors --json
node tools/vibit inspect pitaya-cluster-observability --json
```

Expected: all exit 0 and report false implementation deferrals.

- [ ] **Step 3: Add runtime checks**

Wire a table-driven runtime check that validates source artifacts, manifest references, inspect output, and forbidden direct compatibility markers.

### Task 5: Verification, Commit, Push

**Files:**
- All touched files

- [ ] **Step 1: Run targeted verification**

Run:

```bash
node -c tools/vibit
node tools/vibit inspect pitaya-node-identity --json
node tools/vibit inspect pitaya-cluster-observability --json
node tools/vibit inspect rule runtime.pitaya_aligned_node_identity_boundary_gate
node tools/vibit inspect rule runtime.pitaya_aligned_cluster_observability_source_first_map
```

Expected: all exit 0.

- [ ] **Step 2: Run repository checks**

Run:

```bash
node tools/vibit check runtime --json
node tools/vibit check work --json
node tools/vibit check schemas --json
node tools/vibit check memory --json
node tools/vibit check all --json
git diff --check
```

Expected: all pass.

- [ ] **Step 3: Run Go tests**

Run:

```bash
cd runtime && go test ./...
```

Expected: all packages pass.

- [ ] **Step 4: Scan for committed secrets**

Run:

```bash
rg -n 'ghp_[A-Za-z0-9_]{20,}|github_pat_[A-Za-z0-9_]+' --glob '!node_modules/**' --glob '!.git/**' --glob '!.vibit.local.env'
```

Expected: no output.

- [ ] **Step 5: Commit and push**

Run:

```bash
git status --short --ignored
git add .
git commit -m "Add Pitaya distributed operations source-first maps"
```

Expected: ignored `.vibit.local.env` and `node_modules/` are not staged. Push with the ignored `GITHUB_TOKEN` from `.vibit.local.env` without printing the token.
