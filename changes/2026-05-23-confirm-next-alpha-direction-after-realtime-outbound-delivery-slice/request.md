# Request

## Original Request

```text
那这样吧，现在就是参考那卡马吧，以后再说皮塔亚的事。然后还有一点就是，我们要研发跟测试都是AI原生，也就是用户说了一个需求之后，就AI把事情都帮他办好，这就是我们的这个产品的设计目的，所有的架构跟实现都是为这点服务的。你按这个把计划再调整一下。
```

English summary: adjust the plan so vibit uses Nakama as the current reference, leaves Pitaya for later, and treats AI-native development plus AI-native testing as the product design purpose. When a user states a requirement, AI should handle the work, and all architecture and implementation should serve that purpose.

## Clarified Requirement

Advance `M-147/W-0219 Confirm next alpha direction after realtime outbound delivery slice`.

This is a direction-selection and roadmap posture change only. It must select exactly one bounded next direction, update the reference posture, and open one follow-up work item.

## Selected Direction

Select:

```text
define_agent_native_feature_request_test_workflow
```

as the next prototype-ready alpha direction.

## User-Visible Outcome

Maintainers and agents can see that the next work is not another runtime feature. The next work defines how a user requirement becomes an AI-written spec, acceptance criteria, test plan, tests, implementation, verification, and durable project memory.

## Nakama/Pitaya Alignment

Nakama is now the primary product capability reference.

Pitaya is deferred as a future architecture reference for distributed runtime concerns such as acceptors, sessions, route handlers, remotes/RPC, frontend/backend roles, groups, broadcast, serializers, and service discovery.

Direct Nakama/Pitaya API compatibility remains deferred.

## Non-Goals

- Implementing runtime behavior.
- Adding protocol routes or bridge behavior.
- Adding Protobuf source or generated output.
- Adding startup wiring.
- Adding persistence, migrations, repository changes, PostgreSQL adapters, or dependencies.
- Adding stream subscriptions, chat, groups, broadcast fanout, matchmaking, match runtime, SDKs, hosted deployment, release artifacts, public announcements, or paid promotion.
- Adding Pitaya-style cluster/RPC/frontend-backend/service-discovery behavior.
- Adding direct Nakama/Pitaya API compatibility.

## Acceptance Criteria

- [x] Exactly one next bounded direction is selected.
- [x] Nakama is recorded as the primary product capability reference.
- [x] Pitaya is recorded as a deferred future architecture reference.
- [x] AI-native development and AI-native testing are recorded as the product purpose.
- [x] One follow-up work item is opened.
- [x] Runtime behavior, protocol, generated output, migration, dependency, broad product module, distributed runtime, and direct compatibility deferrals are preserved.
