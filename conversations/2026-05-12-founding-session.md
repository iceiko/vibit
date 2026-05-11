# Conversation: Founding Direction And First Publication

Date: 2026-05-12  
Participants: Maintainer, Agent  
Related changes:

- `changes/2026-05-12-add-conversation-log-standard/`

Related artifacts:

- `CONSTITUTION.md`
- `CONSTITUTION.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `.arch/`
- `docs/module-manifest.md`
- `docs/change-spec.md`
- `conversations/`

## Context

The maintainer explored whether there was value in creating an open-source AI-native server architecture. The initial discussion compared popular game server frameworks and AI-native interpretations. The maintainer then clarified that "AI-native server" should not primarily mean a server that exposes AI gameplay features. It should mean a server architecture designed from first principles for AI coding agents to write, understand, extend, and maintain.

## Maintainer Narrative

The maintainer clarified the core concept:

> 我所谓的AI原生的服务器，是说它是完全由Agent编写的服务器。然后所有的需求都是由Agent继续处理的。这就依赖于所有的Agent能够按照一个既定的标准来完成所有的需求，并且在有新的需求的时候，它能完美地契合到原有的架构中。

The maintainer distinguished this from merely adding AI features:

> 当然，提供AI功能，那是另外一个小的方向。只是说，比如说现在我接触的很多服务器代码，你让AI去看的时候，因为服务器代码的质量以及编程的规范范式的不同，AI有的时候使不上力。我希望有一个全新的服务器架构，能让AI非常顺畅舒适地进行工作。

The maintainer asked for a constitutional document:

> 现在我们就建一个MD文档，作为这个项目的宪法。然后我们会持续维护这个文档。

The maintainer set the bilingual documentation rule:

> 第一，我们所有的文档要用英文撰写，但是要把它翻译成一个中文版本，供人阅读。

The maintainer emphasized naming and product clarity:

> 我们要给这个项目起一个好听的名字，而且要让理解它或者在心里有懵懂想法的人，一眼就能看出来，这是他想要的一个东西。

The maintainer stated the broader thesis:

> AI编程现在风靡到现在，能够推倒重来的，完全按照AI Agent的需求去设计的软件，我觉得才是下一个时代的未来。

The maintainer later selected the product name:

> vibit 这是我给这个产品起的名字。

The maintainer wanted to publish immediately and evolve in public:

> 我们要先把我们现在的东西发布上去，然后我们逐步的去丰富这个产品。我们的协议可以是MIT。

The maintainer then added the project-memory requirement:

> 我作为主开发者，我的每一次与Agent的对话都要进行记录，尤其是我的叙述。Agent的回答可以简洁的记录，这样能够未来能够理清曾经我们怎么走到了这里。

## English Summary Of Maintainer Intent

The maintainer's intent is to create an open-source server framework designed for agent-driven software development from the beginning.

Key points:

- The project is not primarily about AI gameplay features.
- The architecture should make AI coding agents effective by giving them explicit standards, module boundaries, contracts, tests, manifests, and workflows.
- The project should be rebuilt from first principles around agent maintainability rather than retrofitted from conventional server architecture.
- English should be the canonical documentation language, with Simplified Chinese translations for human readers.
- Conversation history is part of project memory and should preserve maintainer narrative with high fidelity.

## Agent Response Summary

The agent recommended defining the project as an agent-native server framework and creating a constitutional document before implementation.

The agent drafted:

- Project constitution
- Bilingual documentation standard
- Naming standard
- Repository agent guide
- README
- MIT license
- Architecture manifests
- Module manifest standard
- Change spec standard

The agent initially proposed `Framewright` as a working name. The maintainer later chose `vibit`; the documents were updated to use `vibit` as the product name while preserving `Framewright` in naming history.

The agent used a temporary GitHub token provided by the maintainer to create and publish:

```text
https://github.com/iceiko/vibit
```

The token itself was not written to repository files or Git remote configuration.

## Decisions

- Product name: `vibit`
- Category: `Agent-Native Server Framework`
- License: MIT
- Repository: `https://github.com/iceiko/vibit`
- Documentation language: English canonical, Simplified Chinese translation required for public-facing docs
- AI-native definition: agent-native maintainability first, AI gameplay features second
- Standards-first approach: constitution, agent guide, manifests, module standard, change standard before implementation
- Conversation logs are project memory and should be preserved under `conversations/`

## Artifacts

Created and published:

- `CONSTITUTION.md`
- `CONSTITUTION.zh-CN.md`
- `AGENTS.md`
- `AGENTS.zh-CN.md`
- `README.md`
- `README.zh-CN.md`
- `LICENSE`
- `.arch/README.md`
- `.arch/README.zh-CN.md`
- `.arch/modules.yaml`
- `.arch/conventions.yaml`
- `docs/module-manifest.md`
- `docs/module-manifest.zh-CN.md`
- `docs/change-spec.md`
- `docs/change-spec.zh-CN.md`
- `changes/_template/`

Created in this change:

- `docs/conversation-log.md`
- `docs/conversation-log.zh-CN.md`
- `conversations/README.md`
- `conversations/README.zh-CN.md`
- `conversations/_template/session.md`
- `conversations/2026-05-12-founding-session.md`

## Open Questions

- Which implementation language should be used for the first CLI prototype?
- Should the first runtime target be game-backend-specific or backend-general?
- How much of the conversation logging process should eventually be automated?
- Should private conversation appendices exist outside the public repository for sensitive context?

## Follow-Up

- Revoke the temporary GitHub token used for publication.
- Add this conversation-log standard to root docs and architecture conventions.
- Start `changes/2026-05-12-bootstrap-vibit-cli/`.
- Build the first CLI prototype for architecture checks and generation.

## Redaction Notes

The maintainer provided a temporary GitHub token during publication. The raw token is intentionally not recorded here and must not be committed to the repository.

Use `[REDACTED_SECRET]` for any future secret-like value that appears in maintainer-agent conversations.
