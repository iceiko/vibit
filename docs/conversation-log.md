# Conversation Log Standard

Status: Draft v0.1  
Last updated: 2026-05-12  
Scope: `conversations/`

This document defines how vibit records maintainer-agent conversations as durable project memory.

Conversation logs do not replace the constitution, change specs, architecture manifests, or issue trackers. They explain the reasoning path behind them.

## 1. Purpose

vibit is being designed through close maintainer-agent collaboration.

The maintainer's narrative is a primary source of product intent. It often contains:

- The problem as felt before it has formal language
- Rejections of misleading interpretations
- Naming intent
- Architectural taste
- Product boundaries
- Historical reasoning behind standards

Future agents should not have to infer this history from final documents alone.

## 2. Location

Conversation logs live under:

```text
conversations/
```

Reusable template:

```text
conversations/_template/session.md
```

Session logs should use:

```text
conversations/YYYY-MM-DD-short-session-id.md
```

## 3. Canonical Language Rule

The project documentation language is English, but conversation logs have a special rule:

- Maintainer statements should be preserved in their original language when possible.
- English summaries should be added for global readability.
- Agent responses may be summarized in English.
- Chinese maintainer statements do not need to be translated line by line, but key decisions and intent should be summarized in English.

This exception exists because the original maintainer wording is part of project memory.

## 4. What To Record

Record:

- Date
- Participants
- Context
- Maintainer narrative, with high fidelity
- Agent response summary
- Decisions made
- Artifacts created or changed
- Open questions
- Follow-up actions

Prefer preserving the maintainer's exact wording when it contains product intent, architectural judgment, or naming rationale.

## 5. What To Summarize

Agent responses may be summarized unless the exact wording is needed for a decision.

Summaries should capture:

- Recommendations
- Alternatives considered
- Warnings or constraints
- Concrete actions taken
- Verification status

Do not preserve long agent output just because it exists.

## 6. What Not To Record

Never commit:

- Access tokens
- API keys
- Passwords
- Private keys
- Session cookies
- One-time codes
- Unrelated personal data
- Private account details not needed for project history

If a secret appears in a conversation, replace it with:

```text
[REDACTED_SECRET]
```

If a private account identifier is not needed, replace it with:

```text
[REDACTED_ACCOUNT]
```

## 7. Relationship To Change Specs

Conversation logs explain how the project got somewhere.

Change specs explain how a specific change is executed.

When a conversation leads to a non-trivial change, link both directions:

- The conversation log should mention the change spec.
- The change spec should mention the conversation that motivated it.

## 8. Required Session Sections

Each session log should include:

```text
# Conversation: <title>

Date:
Participants:
Related changes:
Related artifacts:

## Context

## Maintainer Narrative

## Agent Response Summary

## Decisions

## Artifacts

## Open Questions

## Follow-Up

## Redaction Notes
```

## 9. Agent Rules

Agents should update conversation logs when:

- The maintainer introduces or clarifies product direction.
- The maintainer names or renames a concept.
- The maintainer rejects an interpretation.
- A governance, architecture, or workflow standard is created.
- A meaningful tradeoff is discussed.
- A public artifact is created from a conversation.

For routine implementation work, logs may be brief.

## 10. Verification

Before committing conversation logs, run a secret scan appropriate to the repository.

Initial minimal check:

```bash
rg -n "ghp_|github_pat_|TOKEN|Token:|api[_-]?key|password|secret" .
```

This check is not complete security tooling, but it prevents the most obvious mistakes until dedicated tooling exists.
