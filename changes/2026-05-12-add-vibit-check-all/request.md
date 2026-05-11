# Request

## Original Request

The maintainer approved continuing with the agent's recommended next step and noted that the server implementation language and overall server instance architecture have not yet been selected.

Original maintainer statement:

> 按你的建议继续。我们是不是还没有选服务器的开发语言，也没有定整体的服务器实例架构。不过可以先按照你的节奏来推进，不要太在意我的想法。

## Clarified Requirement

Add an aggregate CLI command:

```text
vibit check all
```

The command should run all currently meaningful repository checks:

- Architecture checks
- Change spec checks
- Module checks

The server implementation language and server instance architecture must remain explicit open decisions. This change should not choose them.

## User-Visible Outcome

A maintainer or agent can run one command to verify the current repository standards:

```bash
node tools/vibit check all
```

## Non-Goals

- Do not choose the server runtime language.
- Do not define the server instance architecture.
- Do not add external CLI dependencies.
- Do not replace individual check commands.

## Unknowns

- How strict `check all` should become once the repository has many modules and many historical change specs.
- Whether future checks should support warnings-as-errors or machine-readable output.

## Acceptance Criteria

- [x] `node tools/vibit check all` exists.
- [x] It runs architecture checks.
- [x] It discovers and checks change specs under `changes/`.
- [x] It discovers and checks modules registered in `.arch/modules.yaml`.
- [x] It exits non-zero if any subcheck fails.
- [x] README and AGENTS mention the command.
- [x] Conversation log records that server language and instance architecture remain open decisions.
