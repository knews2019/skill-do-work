# Skill Anatomy

A guide to how action files are structured, how routing works, and how the pieces fit together.

## Action file structure

Every action file in `actions/` follows a consistent anatomy. The exemplar is `review-work.md`, which demonstrates all sections:

```
---
name: review-work                          ← YAML frontmatter (machine-readable metadata)
description: "Use when..."                 ← Activation hint for agents
---

# Review Work Action                       ← Title

> **Part of the do-work skill.** ...       ← Description blockquote (human-readable summary)

## When to Use                             ← Positive triggers and explicit exclusions
## Philosophy                              ← Why the action exists and design principles
## Input                                   ← What parameters drive behavior ($ARGUMENTS, modes)
## Steps                                   ← Numbered procedural steps (the core workflow)
  ### Step 1: ...
  ### Step 2: ...
## Output Format                           ← What gets produced (report structure, files, output)
## Rules                                   ← Constraints, what NOT to do
## Common Rationalizations                 ← Anti-shortcut table (3-column format)
## Red Flags                               ← Observable failure symptoms
## Verification Checklist                  ← Exit criteria with evidence requirements
```

### Required elements

- **Frontmatter** — `name` and `description` fields. Description starts with "Use when..." for agent activation.
- **Description blockquote** — One sentence starting with `> **Part of the do-work skill.**`
- **Steps** — Numbered `### Step N:` headings with actionable instructions.

### Encouraged elements

- **When to Use** — Helps agents in environments without the SKILL.md routing table.
- **Common Rationalizations** — Guards against agents convincing themselves to skip steps.
- **Red Flags** — Helps reviewers detect problems after the fact.
- **Verification Checklist** — Concrete exit criteria so "done" means provably done.

### Section order

The order matters: **Frontmatter → Philosophy → When to Use → Input → Steps → Output → Rules → Common Rationalizations → Red Flags → Verification Checklist**. Not all sections are required in every file, but the ones present should follow this order.

## Accepted variants

Not all actions follow the flat Steps pattern:

- **Sub-command dispatchers** (`prime.md`, `build-knowledge-base.md`) — Use a Sub-Commands table instead of flat steps. Each sub-command has its own workflow section.
- **Multi-mode actions** (`present-work.md`, `review-work.md`) — Use a Modes table, then separate workflow sections per mode.
- **State-based actions** (`version.md`, `pipeline.md`) — Response sections keyed by input type instead of sequential steps.
- **Companion files** (`work-reference.md`, `deep-explore-reference.md`) — Reference data for another action. Not invoked directly.

## How routing works

`SKILL.md` is the single entry point. It contains:

1. A **priority routing table** (first-match-wins) that maps user input patterns to action files.
2. A **verb reference** table listing all trigger words for each route.
3. **Payload preservation rules** ensuring user content is never lost during routing.

When the user types `do work [something]`, SKILL.md pattern-matches against the table and dispatches to the matching action file. The action file is loaded on-demand — only the active action needs to be in context.

## How crew-members work

Domain-specific rules in `crew-members/` are loaded just-in-time during implementation:

- `general.md` — always loaded during implementation
- `{domain}.md` — loaded when the REQ's `domain` frontmatter matches
- `testing.md` — loaded when `tdd: true` or after 2+ test failures
- `debugging.md` — loaded during remediation and after 2+ test failures

Each crew-member file has a Role Identity section that tells the agent what perspective to adopt. Missing files don't block — the agent proceeds without them.

## How specs work

Templates in `specs/` define output structure and quality standards for common task types (API endpoints, UI components, refactoring, bug fixes). They're:

- Optionally hinted during capture via `suggested_spec` frontmatter
- Automatically loaded during the work action after triage
- Guidance for the builder and reviewer — not binding; REQ requirements take priority

## Slash commands

`.claude/commands/` provides tab-completable shortcuts for Claude Code users. Each command is a thin wrapper that reads the relevant action file and follows it. The SKILL.md routing table remains the canonical dispatch mechanism.

## Cross-referencing

Reference other actions by short name ("the work action", "do work clarify") — never by file path. SKILL.md owns the file-path mappings. This keeps action files portable and avoids breaking references when files move.
