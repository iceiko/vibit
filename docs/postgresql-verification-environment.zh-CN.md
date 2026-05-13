# PostgreSQL Verification Environment Standard 中文版

状态：Draft v0.1  
最后更新：2026-05-13  
范围：用于 live persistence 和 migration verification 的一次性 PostgreSQL environment rules  
说明：本文件是 `docs/postgresql-verification-environment.md` 的简体中文译本。英文版本是权威版本，本译本用于人类阅读、讨论和维护共识。

本标准定义 agents 和 maintainers 如何在不依赖维护者记忆的情况下验证 PostgreSQL-backed behavior。

本标准应与 `docs/postgresql-persistence-boundary.md`、`.arch/runtime.yaml`、`.arch/conventions.yaml`、`ADR-0011`、`ADR-0013` 和 `ADR-0020` 一起使用。

## 1. 目的

Source checks 不等于 live database verification。

`node tools/vibit check migrations` 验证 migration files、ownership traces 和 source conventions。它不能证明 migrations 能在 PostgreSQL 上干净 apply，也不能证明 repository SQL 能针对真实 server 运行，或 transaction behavior 能在真实 database constraints 下成立。

本标准定义这些 live checks 所需的最小显式 environment contract，同时让普通 repository verification 保持快速和本地化。

## 2. Dependency Position

Disposable PostgreSQL verification 默认是 optional。

规则：

- 普通 unit tests 和 `node tools/vibit check all` 不得要求运行中的 PostgreSQL server。
- Live PostgreSQL checks 必须通过显式 environment variables 选择性启用。
- Human 或 agent 可以使用 Docker、Podman、system PostgreSQL 或其他 service manager，但本标准不把任何一种列为 required project dependency。
- 本标准不假设 cloud-hosted PostgreSQL。
- Agents 不得把 credentials 隐藏在 tracked files 中。

## 3. Environment Variables

Live verification environment 由这些 variables 描述：

```text
VIBIT_POSTGRES_TEST_DSN
VIBIT_POSTGRES_TEST_DATABASE
VIBIT_POSTGRES_TEST_CLEANUP
VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE
```

### `VIBIT_POSTGRES_TEST_DSN`

Live PostgreSQL verification 必填。

这个 DSN 指向 agent 被允许在 verification 中修改的一次性 database 或 database server。

规则：

- 该值必须来自显式 local environment input。
- 该值不得保存到 tracked files。
- 如果未设置，live PostgreSQL checks 必须 skip，并记录没有可用 disposable environment。

### `VIBIT_POSTGRES_TEST_DATABASE`

Optional database name。当 DSN 没有完全标识 target database 时，scripts 或 manual setup instructions 可使用它。

推荐默认值：

```text
vibit_test
```

### `VIBIT_POSTGRES_TEST_CLEANUP`

Optional cleanup mode。

允许值：

```text
drop_schema
truncate
keep
```

默认值：

```text
drop_schema
```

规则：

- `drop_schema` 是推荐的一次性模式，因为它能让 database 为干净 rerun 做好准备。
- 当无法 drop schema objects 时，可以使用 `truncate`。
- `keep` 必须记录在 verification output 中，因为它会留下 state 供检查。

### `VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE`

Destructive setup 或 cleanup 必填。

允许值：

```text
0
1
```

默认值：

```text
0
```

规则：

- 除非该值为 `1`，否则 destructive setup 或 cleanup 不得运行。
- Destructive 指 dropping schemas、dropping databases、truncating module-owned tables，或运行会移除 schema objects 的 rollback checks。
- 如果某个 check 需要 destructive behavior，但该值不是 `1`，它必须 skip 或 fail，并给出清晰说明，而不是猜测意图。

## 4. Verification Categories

PostgreSQL verification 分为三类。

### Source Verification

Source verification 默认始终可以安全运行：

```bash
node tools/vibit check migrations
node tools/vibit check postgres-env
```

它验证 repository artifacts 和 standards。它不会打开 database connection。

### Unit Verification

Unit verification 默认可以安全运行：

```bash
cd runtime && go test ./...
cd runtime && go vet ./...
node tools/vibit check runtime
```

Unit tests 可以使用 fake executors 或 fake database handles。它们不得要求 `VIBIT_POSTGRES_TEST_DSN`。

### Live PostgreSQL Verification

Live PostgreSQL verification 是 opt-in。

它可以包括：

```text
针对 VIBIT_POSTGRES_TEST_DSN 的 migration status
针对 VIBIT_POSTGRES_TEST_DSN 的 migration apply
针对 VIBIT_POSTGRES_TEST_DSN 的 repository integration tests
针对 VIBIT_POSTGRES_TEST_DSN 的 transaction runner integration tests
当 destructive verification 被显式允许时，运行 rollback 或 cleanup checks
```

在 live verification commands 存在前，agents 必须记录：

```text
Not verified: live PostgreSQL verification is unavailable because no repository command exists yet.
```

或者，如果 commands 已存在但 `VIBIT_POSTGRES_TEST_DSN` 未设置：

```text
Not verified: live PostgreSQL verification skipped because VIBIT_POSTGRES_TEST_DSN was not set.
```

## 5. Setup Expectations

本标准有意不规定某一种 service manager。

任何 setup path 都可以接受，只要它提供：

- 与第一版 runtime migration 和 repository code 兼容的 PostgreSQL。
- 可以被修改的一次性 database 或 schema。
- 通过 `VIBIT_POSTGRES_TEST_DSN` 提供的 DSN。
- 与 `VIBIT_POSTGRES_TEST_CLEANUP` 匹配的 cleanup behavior。
- 当 cleanup 可能移除 schema 或 data 时，通过 `VIBIT_POSTGRES_TEST_ALLOW_DESTRUCTIVE=1` 明确授予 destructive permission。

Agents 可以记录自己实际使用的 local command，但不得 commit local service credentials、socket paths、passwords 或 tokens。

## 6. Cleanup Rules

Agents 必须让 environment 处在已知状态。

规则：

- 优先使用 isolated databases 或 isolated schemas。
- 成功 live verification run 后，优先清理创建的 schema objects。
- 如果有意跳过 cleanup，记录原因。
- 如果 cleanup 失败，把它记录为 verification risk。
- 永远不要针对未明确提供为 disposable verification 的 DSN 运行 destructive cleanup。

## 7. Recording Verification

任何触及 PostgreSQL migrations、repositories、transaction boundaries 或 persistent runtime composition 的 change，都必须记录以下之一：

```text
Verified: live PostgreSQL verification ran with VIBIT_POSTGRES_TEST_DSN set.
Not verified: live PostgreSQL verification skipped because VIBIT_POSTGRES_TEST_DSN was not set.
Not verified: live PostgreSQL verification command is not implemented yet.
Not applicable: change did not touch PostgreSQL-backed behavior.
```

当 live verification 运行时，记录：

- 运行过的 commands。
- cleanup 是否运行。
- 是否允许 destructive cleanup。
- 任何 skipped integration tests。

不要记录 DSN 本身。

## 8. Current Repository Tooling

当前 static environment-standard check 是：

```bash
node tools/vibit check postgres-env
```

该检查验证 disposable PostgreSQL verification standard、runtime manifest references 和 guidance artifacts 是否存在。它不会连接 PostgreSQL。

Live migration 和 repository integration commands 仍属于未来工作。
