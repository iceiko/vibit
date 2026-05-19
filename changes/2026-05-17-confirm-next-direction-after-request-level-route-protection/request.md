# Request

## Original Request

The maintainer asked:

```text
选择 wire_runtime_authentication_startup_composition，继续推进十步，注意，重点参考。Nakama 和 Pitaya。这些已有的游戏Server，他们对现在游戏的需求理解得非常透彻，我们要重点参考。
```

## Clarified Requirement

Close the blocked `M-043/W-0115` next-direction confirmation gate and select:

```text
wire_runtime_authentication_startup_composition
```

The selected direction must create a bounded startup-composition gate before implementation and must keep Nakama and Pitaya as active reference baselines without adopting direct API compatibility.

## User-Visible Outcome

The work queue moves from the blocked request-level route protection direction gate into a startup composition milestone, starting with a gate-only work item.

## Non-Goals

- Do not add session persistence.
- Do not add WebSocket handshake authentication.
- Do not add logout, refresh, cleanup, or token rotation behavior.
- Do not change repository interfaces, PostgreSQL adapters, migrations, generated files, or dependencies.
- Do not add authentication command Protobuf messages beyond the existing wrapper.
- Do not expand core game backend modules in this direction-confirmation change.

## Acceptance Criteria

- [x] `W-0115` is marked completed with `selected_direction: wire_runtime_authentication_startup_composition`.
- [x] `M-043` is marked completed.
- [x] `M-044` is created as the startup composition gate milestone.
- [x] `W-0116` is created as the first startup composition gate work item.
- [x] Nakama and Pitaya reference implications are recorded.
- [x] Verification is recorded.
