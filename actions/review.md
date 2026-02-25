# Review Action

> **Part of the do-work skill.** Invoked automatically after work completes or manually when the user requests a code review. Evaluates implementation quality against the original requirements.

A post-work quality gate that reads the actual code changes and compares them against the REQ and UR to catch what automated tests miss: requirements gaps, code quality issues, scope creep, and missing edge cases.

## Philosophy

- **Review is not testing.** Tests verify behavior; review evaluates quality, completeness, and intent alignment.
- **Traceability makes this powerful.** You have the original input (UR), the structured requirements (REQ), the triage/plan/exploration history, and the actual diff — use all of it.
- **Actionable output.** Don't just report problems — create follow-up REQs for issues worth fixing.
- **Proportional effort.** A Route A config change gets a quick scan. A Route C multi-file feature gets a thorough review.

## Two Modes

| Mode | Trigger | REQ location | How to get the diff |
|------|---------|-------------|---------------------|
| **Pipeline** | Auto-triggered by the work action after testing passes | `do-work/working/` | `git diff` (uncommitted changes) or read the files listed in the Implementation Summary |
| **Standalone** | User invokes manually: `do work review`, `do work code review`, `do work review REQ-005` | `do-work/archive/` or `do-work/archive/UR-NNN/` | `git show <commit>` using the `commit` frontmatter field |

Both modes follow the same workflow. The only difference is where the REQ lives and how you obtain the diff.

## Workflow

### Step 1: Find the Target

**Pipeline mode:** The work action hands you the REQ file path (in `do-work/working/`). Skip to Step 2.

**Standalone mode:**
1. **If user specifies a REQ** (e.g., "review REQ-005"): Find it in `do-work/archive/` or `do-work/archive/UR-NNN/`
2. **If user specifies a UR** (e.g., "review UR-003"): Find all completed REQs under that UR and review each
3. **If no target specified**: Find the most recently completed REQ — check `do-work/archive/` for the highest REQ number with `status: completed`

If the target REQ has no `commit` field (standalone mode) or no implementation changes (pipeline mode), report that there's nothing to review and exit.

### Step 2: Read the REQ

Read the full REQ file. Extract:
- **What was requested** — the What/Detailed Requirements sections
- **Triage decision** — the route and reasoning
- **Plan** — what was planned (if Route B/C)
- **Implementation Summary** — what the builder says it did
- **Testing results** — what tests exist and pass

### Step 3: Read the Original Input

Find the UR via the REQ's `user_request` frontmatter field. Read `do-work/user-requests/UR-NNN/input.md` (or `do-work/archive/UR-NNN/input.md` in standalone mode). This is the source of truth for what the user actually wanted.

If the REQ is a legacy file without `user_request`, use whatever context is available (the REQ content itself, any `context_ref` file).

### Step 4: Get the Diff

**Pipeline mode:** Run `git diff` to see uncommitted changes, or read the files the Implementation Summary lists as created/modified. If the working tree is clean (implementation agent already staged), use `git diff --staged`.

**Standalone mode:** Run `git show <commit>` using the hash from the REQ's `commit` frontmatter. This gives you the full diff of what was committed.

Read the diff carefully. For large diffs, focus on:
- New files created (read them fully)
- Modified files (understand what changed and why)
- Deleted files (was deletion justified?)

### Step 5: Evaluate

Score the implementation across these dimensions:

**Requirements Compliance (0-100%)**
- Does the code deliver every requirement from the REQ?
- Are specific values, constraints, and conditions implemented correctly?
- Are edge cases from the requirements handled?
- Compare the Implementation Summary claims against the actual diff

**Code Quality (0-100%)**
- Does it follow existing project patterns and conventions?
- Naming clarity, readability, maintainability
- Appropriate error handling for the context
- No obvious bugs or logic errors in the diff

**Test Adequacy (0-100%)**
- Are there tests for the new/changed behavior?
- Do tests cover the important paths (not just the happy path)?
- Are test assertions meaningful (not just "doesn't throw")?
- If no tests exist and the project has no test infrastructure, score N/A

**Scope Discipline (0-100%)**
- Did the implementation stay focused on the request?
- Any unnecessary refactoring, feature additions, or style changes?
- Files touched that didn't need touching?

**Risk Assessment (Critical / Low / None)**
- Security concerns (injection, auth bypass, data exposure)
- Performance risks (N+1 queries, unbounded loops, memory leaks)
- Data integrity risks (race conditions, missing validation at boundaries)

### Scoring Guidelines

**90-100%**: Ship-ready. Clean implementation, good tests, on-spec.
**75-89%**: Minor issues. Style nits, a missing edge case test, small scope drift. Not worth blocking.
**50-74%**: Needs attention. Missed requirements, inadequate tests, or code quality concerns.
**Below 50%**: Significant problems. Major requirements missed, no tests for new behavior, or risky code.

### Step 6: Report

**Pipeline mode:** Report to the work action orchestrator (which reports to the user).
**Standalone mode:** Report directly to the user.

Format:

```
## Code Review: REQ-NNN

**Overall: [X]%** | Route [A/B/C] | [commit hash or "uncommitted"]

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 95% | All requirements implemented |
| Code Quality | 85% | Clean, follows patterns |
| Test Adequacy | 70% | Missing edge case coverage |
| Scope | 100% | Focused, no drift |
| Risk | None | — |

### Findings

**Important:**
- [Specific finding with file:line reference]

**Minor:**
- [Style nit or suggestion]

### Follow-up REQs Created
- REQ-025: "Add edge case tests for dark mode toggle" (addendum_to: REQ-003)
```

### Step 7: Create Follow-up REQs

For each **Critical** or **Important** finding, create a follow-up REQ:

```markdown
---
id: REQ-NNN
title: "Review fix: [brief description]"
status: pending
created_at: [timestamp]
user_request: [same UR as the reviewed REQ]
addendum_to: [reviewed REQ id]
review_generated: true
---

# Review Fix: [Brief Description]

## What
[Describe the issue found and the fix needed]

## Context
Found during code review of [REQ-id]. [1 sentence on what the review found.]

## Requirements
- [Specific fix needed]
```

Follow-up REQs go in `do-work/` (the queue). In pipeline mode, the work loop picks them up on the next iteration. In standalone mode, they wait for the user to run `do work run`.

**Don't create follow-ups for minor issues.** Minor findings go in the report only. The threshold: would a senior engineer request changes on this in a PR review, or just leave a comment?

### Append to REQ File

After generating the report, append a Review section to the REQ file:

```markdown
## Review

**Overall: [X]%** | [timestamp]

| Dimension | Score |
|-----------|-------|
| Requirements | X% |
| Code Quality | X% |
| Test Adequacy | X% |
| Scope | X% |
| Risk | [level] |

**Findings:** [count] important, [count] minor
**Follow-ups created:** [REQ-NNN, REQ-NNN] or "None"

*Reviewed by review action*
```

In standalone mode, this is an exception to the archive immutability rule — review annotations are post-work metadata, not content changes.

## Calibrating Review Depth

Match effort to complexity:

| Route | Review depth |
|-------|-------------|
| **A** (Simple) | Quick scan. Check the change is correct and focused. 1-2 minutes of review. Skip dimensions that don't apply (e.g., Test Adequacy for a copy change). |
| **B** (Medium) | Standard review. Check all dimensions. Focus on whether the builder found and followed the right patterns. |
| **C** (Complex) | Thorough review. Compare against the plan. Check architectural decisions. Verify cross-cutting concerns. Read new files fully. |

## What NOT to Do

- Don't re-implement — you're reviewing, not building
- Don't review your own review's follow-up REQs more strictly than the original work — avoid infinite loops of diminishing-return fixes
- Don't block on minor issues — report them but keep moving
- Don't invent requirements — review against what the REQ says, not what you think it should say
- Don't penalize the absence of things the project doesn't have (no test infrastructure = don't fail on test adequacy)
- Don't duplicate what testing already covers — if tests pass and cover a behavior, don't re-verify that behavior manually
