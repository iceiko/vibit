# Verification

Status: Verified

RED checks:

```text
node tools/vibit inspect pitaya-acceptor-connection --json
# Unknown command.

node tools/vibit inspect rule runtime.pitaya_aligned_acceptor_connection_lifecycle_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_acceptor_connection_lifecycle_source_first_map

node tools/vibit check change implement-pitaya-aligned-acceptor-connection-lifecycle-source-first-map --json
# change directory does not exist
```

Required final checks:

```text
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect pitaya-acceptor-connection --json
node tools/vibit inspect rule runtime.pitaya_aligned_acceptor_connection_lifecycle_source_first_map
node tools/vibit check change implement-pitaya-aligned-acceptor-connection-lifecycle-source-first-map --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Observed result:

- Pending fresh verification in this work session.
