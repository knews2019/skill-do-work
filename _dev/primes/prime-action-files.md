# Prime: Action Files

> Read this before adding or modifying any action file under `skills/do-work*/actions/`.
> Prime files are low noise, high value: they point at the code that is the source of
> truth, carry no volatile metrics, and every new lesson in this domain gets linked from
> here so future sessions find it. This file is maintainer-side (`_dev/` is
> export-ignored) — nothing shipped may cite it.

## Every New Action Must Justify Its Package

Before an action file is added, state — in its description blockquote or an accompanying ADR — why it belongs in core, board, knowledge, or toolbox: what package machinery it needs or which existing action it completes. Reviewers reject additions without this justification. Every new action also updates the owning `skills/do-work*/SKILL.md`, help, and any next-step surface; router budgets are enforced by `_dev/tests/contract-regressions.sh`.

## Template

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

[Include only if earned — see below. Constraints specific to this action, not restated engineering hygiene.]

## Common Rationalizations

[Include only if earned — see below.]

| If you're thinking...              | STOP. Instead...     | Because...               |
| ---------------------------------- | -------------------- | ------------------------ |
| [Shortcut the agent might attempt] | [What to do instead] | [Why the shortcut fails] |

## Red Flags

[Include only if earned — see below.]

- [Observable symptom that something went wrong — helps reviewers detect problems after the fact]

## Verification Checklist

[Include only if earned — see below.]

- [ ] [Concrete exit criterion with evidence requirement]
```

**Required:** Description blockquote, Steps (numbered). **Common:** Input, Output Format, When to Use.

**Earned, not mandatory: Rules, Common Rationalizations, Red Flags, Verification Checklist.** Add one only when the file has something a capable model would otherwise get wrong — do-work machinery (a queue/pipeline mechanic, a frontmatter or schema contract) or a hard-won failure mode with a traceable origin (a real REQ or incident this stops from recurring). "This is generic engineering advice a capable model already follows" is an explicit *non*-reason — true or not, it doesn't earn a section.

**The test, not a vibe:** before adding a Common Rationalizations row, ask *can I name the specific failure this row prevents, and where it happened?* No → don't add the row. If every row in a table fails that test, omit the whole section — a generic table is worse than no table: it teaches the reader the section is decorative, so they stop reading the ones that aren't. Apply the same test to Rules and Red Flags — specific to this action, not restated hygiene ("write tests," "don't skip validation"). When a file has nothing that passes, omit the section entirely; don't ship it empty or generic to satisfy the template.

**State intent, not a directive rule, when a capable model can infer the rest.** "Report drift, don't fix it inline" gives the model this action's boundary in one line — a five-line Rules section re-deriving why inline fixes are bad adds nothing a capable model didn't already know.

`_dev/tests/contract-regressions.sh` locks in the Common Rationalizations rule: a new action file's table must contain at least one do-work-specific noun (illustrative, not exhaustive, per Closed Enumerations Go Stale in `prime-shell-commands.md` — e.g. REQ, UR, queue, frontmatter, pipeline, archive) or the suite fails, naming the file and the fix.

**Section order when present:** Philosophy → When to Use → Input → Steps → Output → Rules → Common Rationalizations → Red Flags → Verification Checklist.

## Accepted Variants

- **Sub-command dispatchers** (`prime.md`, `bkb.md`) — Use a Sub-Commands table instead of flat steps. Each sub-command has its own workflow section.
- **Multi-mode actions** (`present-work.md`, `review-work.md`, `tutorial.md`) — Use a Modes table, then separate workflow sections per mode. A single `Step 1: Mode Selection` dispatcher at the top is acceptable.
- **State-based actions** (`version.md`, `pipeline.md`) — Response sections keyed by input type instead of sequential steps.
- **Checklist-based diagnostics** (`forensics.md`) — Use a `## Checks` section with independently-runnable items instead of ordered `## Steps`. Each check is a diagnostic probe, not a sequential step.

## Cross-Referencing

Cross-reference same-package actions by their local path (for example `actions/work.md`); cross a package boundary with an explicit sibling path such as `../do-work-knowledge/actions/bkb.md`. Shipped files must never cite this repo's own `CLAUDE.md` or `AGENTS.md` — both are export-ignored maintainer instructions, so the citation dangles in every consumer install. `_dev/tests/contract-regressions.sh` enforces this across the shipped `skills/` tree.

## Descriptions Are Triggers, Not Summaries

A skill or action description is loaded whether or not the thing gets used, so it should carry the words that tell the router when to pull it in — not an explanation so complete the file never gets read.

**Bad:** "Monitors a pull request through review and CI, rebasing as needed and addressing reviewer comments until checks pass."

**Good:** "Use when the user asks to monitor, watch, or babysit a PR."

## Lessons

- [REQ-193: key closed-state authority to the archived-fallback condition and return review-generated completions in place](../../do-work/archive/UR-043/REQ-193-keep-archived-urs-closed-during-review.md#lessons-learned)
- [REQ-194: retain canonical structured detector evidence and test the source seam directly](../../do-work/archive/UR-043/REQ-194-forward-stray-reqs-through-forensics.md#lessons-learned)
- [REQ-189: shared instruction contracts must inherit upstream token grammars and align prescribed shell publication timing](../../do-work/archive/REQ-189-canonical-ai-report-and-shared-evidence-contract.md#lessons-learned)
- [REQ-190: delete obsolete action modes before rebuilding dispatch, and publish immutable outputs before mutable canonical files](../../do-work/archive/REQ-190-reduce-present-work-to-portfolio-only.md#lessons-learned)
- [REQ-191: recover useful source contracts from history without restoring unsafe wrappers, and keep shared mechanics canonical](../../do-work/archive/REQ-191-extract-explicit-present-video-action.md#lessons-learned)
- [REQ-192: keep presentation aliases exact, make guardrail callers condition-based, and test live command ownership at canonical seams](../../do-work/archive/REQ-192-migrate-presentation-routing-docs-and-contracts.md#lessons-learned)
- [REQ-197: test canonical-contract inheritance with active ordered directives and adversarial negations, not duplicated examples](../../do-work/archive/REQ-197-normalize-completed-work-presentation-target-ids.md#lessons-learned)
- [REQ-198: keep optional image batches private until status-backed success and replay the prescribed artifact shape](../../do-work/archive/REQ-198-publish-generated-directory-only-after-image-success.md#lessons-learned)
