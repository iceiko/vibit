# Verification

Status: Verified

RED checks:

```text
node tools/vibit inspect pitaya-serializer-forwarding --json
# Unknown command.

node tools/vibit inspect rule runtime.pitaya_aligned_serializer_message_forwarding_source_first_map
# Unknown rule_id: runtime.pitaya_aligned_serializer_message_forwarding_source_first_map

node tools/vibit check change implement-pitaya-aligned-serializer-message-forwarding-source-first-map --json
# change directory does not exist
```

Required final checks:

```text
node -c tools/vibit
node tools/vibit inspect next --json
node tools/vibit inspect pitaya-serializer-forwarding --json
node tools/vibit inspect rule runtime.pitaya_aligned_serializer_message_forwarding_source_first_map
node tools/vibit check change implement-pitaya-aligned-serializer-message-forwarding-source-first-map --json
node tools/vibit check work --json
node tools/vibit check runtime --json
node tools/vibit check memory --json
node tools/vibit check schemas --json
node tools/vibit check all --json
git diff --check
```

Observed result:

- `node -c tools/vibit` - passed.
- `node tools/vibit inspect next --json` - passed; reports W-0264 as next-ready.
- `node tools/vibit inspect pitaya-serializer-forwarding --json` - passed.
- `node tools/vibit inspect rule runtime.pitaya_aligned_serializer_message_forwarding_source_first_map` - passed.
- `node tools/vibit check change implement-pitaya-aligned-serializer-message-forwarding-source-first-map --json` - passed with 13 passed, 0 warnings, 0 failures.
- `node tools/vibit check work --json` - passed with 1578 passed, 0 warnings, 0 failures.
- `node tools/vibit check runtime --json` - passed with 22420 passed, 1 existing warning, 0 failures.
- `node tools/vibit check memory --json` - passed with 4484 passed, 0 warnings, 0 failures.
- `node tools/vibit check schemas --json` - passed with 4876 passed, 0 warnings, 0 failures.
- `node tools/vibit check all --json` - passed with 315 subchecks passed, 1 existing warning, 0 failures.
- `git diff --check` - passed.
