# Do-Work Skill Project

A task queue skill for agentic coding tools. Platform-agnostic — works with any agent that can read/write files and run shell commands.

## Project Structure

```
SKILL.md              # Entry point — routing logic, action dispatch
next-steps.md         # Per-action next-step suggestions (referenced by SKILL.md)
README.md             # Installation + quick usage guide
actions/              # Action files (each is a standalone prompt)
  capture.md          # Capture new requests → UR folders + REQ files
  work.md             # Process the queue — triage, plan, build, test, review
  work-reference.md   # Orchestrator checklist, error handling, progress template
  clarify.md          # Batch-review pending questions from completed work
  verify-requests.md  # Quality-check captured REQs against original input
  review-work.md      # Post-work code review + acceptance testing
  code-review.md      # Standalone codebase review — consistency, patterns, security, performance, test coverage
  ui-review.md        # Read-only UI quality validation against design best practices
  present-work.md     # Client-facing deliverables (briefs, videos, diagrams)
  cleanup.md          # Archive consolidation
  commit.md           # Atomic git commits traced to REQs
  inspect.md          # Explain uncommitted changes — what, why, and readiness (read-only)
  version.md          # Version reporting + update checks (current version lives here)
  quick-wins.md       # Scan for refactoring opportunities and low-hanging tests
  scan-ideas.md       # Generate ideas for what to build, improve, or explore next
  deep-explore.md     # Multi-round structured exploration — diverge/converge dialogue, vision docs
  deep-explore-reference.md # Companion: persona prompts, rubrics, state schema, error handling
  install-ui-design.md # Install the frontend-design skill for UI work
  install-bowser.md   # Install Playwright CLI + Bowser skill for browser automation
  forensics.md        # Pipeline diagnostics — stuck work, hollow completions, orphaned URs
  prime.md             # Prime file management — create and audit AI context documents
  pipeline.md          # Full end-to-end orchestration (investigate → capture → verify → run → review)
  build-knowledge-base.md # LLM knowledge base — init, triage, ingest, query, lint, and more
  bkb-reference.md    # Companion: seed file templates, agent crew definitions, KB schema content
  tutorial.md          # Interactive tutorials — quick start, concepts, recipes, guided tour
  sample-archived-req.md # Example of a fully processed REQ file (reference only)
specs/                # Reusable specification templates for common task types
  README.md           # What specs are, how to use them, how to create new ones
  api-endpoint.md     # Spec template for building API endpoints
  ui-component.md     # Spec template for frontend UI components
  refactor.md         # Spec template for refactoring tasks
  bug-fix.md          # Spec template for bug fixes
crew-members/         # Domain-specific rules loaded by work action
hooks/                # Optional hook scripts (platform-specific, installable)
  hooks.json          # Combined hook config for Claude Code (SessionStart + Stop)
  session-start.sh    # Claude Code SessionStart hook — injects status line
  pipeline-guard.sh   # Claude Code Stop hook — prevents stopping mid-pipeline
docs/                 # Per-action user guides (capture-guide.md, work-guide.md, etc.)
AGENTS.md             # Stub — redirects to CLAUDE.md
CHANGELOG.md          # Release notes (newest on top)
```

## Before Every Commit

1. **Bump the version** in `actions/version.md` (line starting with `**Current version**:`). Use semver — patch for fixes, minor for features, major for breaking changes. When in doubt, patch. **Verify the new version number is strictly greater than the first existing entry in `CHANGELOG.md`** — duplicate version numbers have occurred before.

2. **Add a changelog entry** at the top of `CHANGELOG.md` (below the header). **Verify the codename is not already used** by an earlier entry.

```markdown
## X.Y.Z — The [Fun Two-Word Name] (YYYY-MM-DD)

[1-2 casual sentences — what changed and why it matters.]

- [Bullet points for specifics]
```

Keep it brief, newest on top, lead with value not implementation. Every version gets an entry.

## Action File Conventions

Action files follow a consistent structure. When adding or modifying actions, use this template:

```markdown
# [Action Name] Action

> **Part of the do-work skill.** [1 sentence: what it does and when it's invoked.]

[Optional: read-only flag, philosophy, or key principles — 1-2 paragraphs max]

## When to Use

**Use when:** [2-4 bullets — positive triggers]
**Do NOT use when:** [2-3 bullets — explicit exclusions, with redirect to correct action]

## Input

[What parameters drive behavior: $ARGUMENTS, target REQ/UR, modes]

## Steps

### Step 1: [First action]
### Step 2: [...]
### Step N: [Final action]

## Output Format

[What gets produced — report structure, file changes, or user-facing output]

## Rules

[Constraints, common mistakes, what NOT to do]

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| [Shortcut the agent might attempt] | [What to do instead] | [Why the shortcut fails] |

## Red Flags

- [Observable symptom that something went wrong — helps reviewers detect problems after the fact]

## Verification Checklist

- [ ] [Concrete exit criterion with evidence requirement]
```

**Required elements:** Description blockquote, Steps (numbered). **Common elements:** Input, Output Format, Rules, When to Use. **Encouraged elements:** Common Rationalizations, Red Flags, Verification Checklist. **Section order matters:** always Philosophy → When to Use → Input → Steps → Output → Rules → Common Rationalizations → Red Flags → Verification Checklist.

**Accepted variants:**
- **Sub-command dispatchers** (`prime.md`, `build-knowledge-base.md`) — Use a Sub-Commands table instead of flat steps. Each sub-command has its own workflow section.
- **Multi-mode actions** (`present-work.md`, `review-work.md`) — Use a Modes table, then separate workflow sections per mode.
- **State-based actions** (`version.md`, `pipeline.md`) — Response sections keyed by input type instead of sequential steps.

Cross-reference other actions by short name (e.g., "the work action", "do work clarify") — not by file path. SKILL.md owns the file-path mappings.

## Agent Rules

Domain-specific rules live in `crew-members/[domain].md`. Each file has a `JIT_CONTEXT` comment documenting when it loads. Loading behavior:

- `general.md` — always loaded during implementation (Step 6)
- `[domain].md` — loaded when the REQ's `domain` frontmatter matches and the file exists
- `testing.md` — loaded when `tdd: true` or `domain: testing`, and alongside debugging.md after 2+ test failures
- `debugging.md` — loaded during remediation (review fail → retry) and after 2+ test failures
- If a rules file is missing, proceed without it — never block on a missing rules file

## Queue Path Convention

Pending REQ files live in `do-work/queue/`. When referencing the queue in action files, always use `do-work/queue/` — not `do-work/` root.

## Agent Compatibility

Action files must work with **any** agentic coding tool:

- Use generalized language ("spawn a subagent", "use your environment's ask-user prompt") — no tool-specific APIs in action files.
- Each action file should work as a standalone prompt pasted into a basic chat interface.
- Design for the floor: the simplest agent that can read/write files and run shell commands must be able to follow the instructions. Subagents and parallel execution are nice-to-haves.
