# Review Work Action

> **Part of the do-work skill.** Invoked automatically after work completes or manually when the user requests a review. Evaluates whether the work actually delivers what was requested — through requirements checking, code review, acceptance testing, and additional testing recommendations.

A post-work quality gate with three jobs: (1) confirm the implementation matches the requirements, (2) verify the code is solid, and (3) actually test that the thing works. Creates follow-up REQs for anything that needs fixing.

## Philosophy

- **Did we build what was asked?** Requirements check comes first. Everything else is secondary if the wrong thing got built.
- **Does it actually work?** Reading diffs catches logic errors. Running the code catches everything else. Do both.
- **Traceability makes this powerful.** You have the original input (UR), the structured requirements (REQ), the triage/plan/exploration history, and the actual diff — use all of it.
- **Actionable output.** Don't just report problems — create follow-up REQs for issues worth fixing.
- **Proportional effort.** A Route A config change gets a quick scan. A Route C multi-file feature gets a thorough review.
- **Suggest what's next.** After checking what you can, tell the user what else should be tested — manual checks, edge cases, integration scenarios, things only a human can verify.

## Two Modes

| Mode | Trigger | REQ location | How to get the diff |
|------|---------|-------------|---------------------|
| **Pipeline** | Auto-triggered by the work action after testing passes | `do-work/working/` | `git diff` (uncommitted changes) or read the files listed in the Implementation Summary |
| **Standalone** | User invokes manually: `do work review`, `do work review work`, `do work review REQ-005` | `do-work/archive/` or `do-work/archive/UR-NNN/` | `git show <commit>` using the `commit` frontmatter field |

Both modes follow the same workflow. The only difference is where the REQ lives and how you obtain the diff.

## Workflow

### Step 1: Find the Target

**Pipeline mode:** The work action hands you the REQ file path (in `do-work/working/`). Skip to Step 2.

**Standalone mode:**
1. **If user specifies a REQ** (e.g., "review REQ-005"): Find it in `do-work/archive/` or `do-work/archive/UR-NNN/`
2. **If user specifies a UR** (e.g., "review UR-003"): Find all completed REQs under that UR and review each
3. **If no target specified**: Find the most recently completed REQ — check both `do-work/archive/` (root) and all `do-work/archive/UR-NNN/` subdirectories for the highest REQ number with `status: completed`

If the target REQ has no `commit` field (standalone mode) or no implementation changes (pipeline mode), report that there's nothing to review and exit.

### Step 2: Read the REQ

Read the full REQ file. Extract:
- **What was requested** — the What/Detailed Requirements sections
- **Builder Guidance** — certainty level (Firm vs Exploratory), scope cues, implementation hints. Use this to calibrate expectations: Exploratory requests get more latitude on interpretation; Firm requirements must match exactly.
- **Triage decision** — the route and reasoning
- **Plan** — what was planned (if Route C)
- **Implementation Summary** — what the builder says it did
- **Testing results** — what tests exist and pass

### Step 3: Read the Original Input

Find the UR via the REQ's `user_request` frontmatter field. Read `do-work/user-requests/UR-NNN/input.md`. If not found there (UR already archived), check `do-work/archive/UR-NNN/input.md`. This is the source of truth for what the user actually wanted.

If the REQ is a legacy file without `user_request`, use whatever context is available (the REQ content itself, any `context_ref` file).

### Step 4: Get the Diff

**Pipeline mode:** Run `git diff` to see uncommitted changes, or read the files the Implementation Summary lists as created/modified. If the working tree is clean (implementation agent already staged), use `git diff --staged`.

**Standalone mode:** Run `git show <commit>` using the hash from the REQ's `commit` frontmatter. This gives you the full diff of what was committed.

Read the diff carefully. For large diffs, focus on:
- New files created (read them fully)
- Modified files (understand what changed and why)
- Deleted files (was deletion justified?)

### Step 5: Requirements Check

Walk through every requirement in the REQ and the original UR input. For each one, determine: **was it built?**

This is not a code quality check — it's a checklist. Go requirement by requirement:

1. **Extract all requirements** from the REQ's What/Detailed Requirements sections AND from the UR's original input. Include explicit requirements, implicit UX expectations, constraints, and edge cases the user mentioned.
2. **For each requirement**, check the diff and implementation:
   - **Delivered**: The requirement is implemented and visible in the diff
   - **Partially delivered**: Some aspects implemented, others missing
   - **Not delivered**: No evidence of implementation in the diff
   - **N/A**: Requirement doesn't apply (e.g., deferred by Open Questions)
3. **Compare against Implementation Summary claims.** If the builder says it did something, verify it actually appears in the diff.

Score: **Requirements Compliance (0-100%)** — percentage of requirements that are fully delivered.

### Step 6: Code Review

Evaluate the implementation quality by reading the diff:

**Code Quality (0-100%)**
- Does it follow existing project patterns and conventions?
- Naming clarity, readability, maintainability
- Appropriate error handling for the context
- No obvious bugs or logic errors in the diff
- Diff hygiene — no debug artifacts, console.log/print statements, or temporary files left behind. **Protect lessons learned** — comments that document *why* something was done, what was tried and didn't work, or architectural reasoning are valuable and should stay. Strip noise, keep knowledge.

**Test Adequacy (0-100%)**
- Are there tests for the new/changed behavior?
- Do tests cover the important paths (not just the happy path)?
- Are test assertions meaningful (not just "doesn't throw")?
- If no tests exist and the project has no test infrastructure, score N/A — this dimension is excluded from the overall average (don't count it as 0%)

**Scope Discipline (0-100%)**
- Did the implementation stay focused on the request?
- Any unnecessary refactoring, feature additions, or style changes?
- Files touched that didn't need touching?

**Risk Assessment (Critical / Low / None)**
- Security concerns (injection, auth bypass, data exposure)
- Performance risks (N+1 queries, unbounded loops, memory leaks)
- Data integrity risks (race conditions, missing validation at boundaries)
- Regression risk — identify callers/dependents of changed code, flag interfaces whose contract changed, note shared utilities that other features rely on

### Step 7: Acceptance Testing

Actually verify the implementation works. Reading diffs catches logic errors; running code catches everything else.

**What to do:**

1. **Run the test suite** — if tests weren't already run by the work pipeline (pipeline mode should have run them in Step 6.5), run them now. Target tests related to changed code first, then broader tests if fast enough.
2. **Try the feature** — if the change produces observable behavior (UI, CLI output, API response, file output), verify it works end-to-end:
   - Run the app/server/tool if applicable
   - Exercise the specific feature that was built
   - Try the happy path first, then obvious edge cases
   - For CLI tools: run the command with expected inputs
   - For APIs: make a test request
   - For UI: describe what you'd check (the reviewer may not have a browser, but should note what to verify)
3. **Verify the fix** — for bug fixes, confirm the original bug no longer reproduces
4. **Check for regressions** — based on the risks identified in code review:
   - Run tests for adjacent/dependent features, not just the changed code
   - If the change modifies a shared utility, API, or interface, exercise the other consumers
   - For bug fixes, verify the fix doesn't break the feature's other behaviors
   - Try the most obvious related flow — if you changed checkout, make sure the cart still works

**What NOT to do:**
- Don't build an exhaustive test harness — this is a quick smoke test, not QA
- Don't test things unrelated to the change
- Don't skip this step just because unit tests pass — unit tests and acceptance testing catch different things

**If you can't run the code** (e.g., requires external services, credentials, or hardware you don't have), note what you couldn't test and include it in the "Suggested Additional Testing" section.

Score: **Acceptance (Pass / Partial / Fail / Untested)**
- **Pass**: Feature works end-to-end as specified
- **Partial**: Core functionality works but edge cases or secondary behaviors don't
- **Fail**: Feature doesn't work as specified
- **Untested**: Couldn't run the code (note why)

### Step 8: Suggest Additional Testing

After completing your review and acceptance testing, recommend what else should be checked. This is where you flag things only a human can verify, things that need specific environments, or edge cases worth exploring.

**Categories to consider:**

- **Manual verification needed** — UI appearance, UX flow, accessibility, things that need eyes on a screen
- **Environment-specific testing** — different browsers, mobile, OS-specific behavior, production-like data
- **Integration testing** — third-party services, APIs, database migrations, auth flows
- **Edge cases worth exploring** — unusual inputs, boundary conditions, concurrent usage, error recovery
- **Performance testing** — if the change could affect performance (large datasets, frequent operations, network calls)
- **Security review** — if the change handles user input, auth, or sensitive data
- **Regression scenarios** — adjacent flows that could break: shared utilities, upstream/downstream consumers, features that depend on the same data or state

Only include categories that are relevant to this specific change. A copy change doesn't need security review suggestions. A new auth flow does.

### Scoring Guidelines

**Computing the overall score:** Average the percentage dimensions — Requirements, Code Quality, Test Adequacy (skip if N/A), and Scope Discipline. Then apply qualitative modifiers: Risk = Critical caps the overall at 60% max; Acceptance = Fail caps at 50% max; Acceptance = Partial applies a 10-point penalty.

**90-100%**: Ship-ready. Clean implementation, good tests, on-spec, acceptance passes.
**75-89%**: Minor issues. Style nits, a missing edge case test, small scope drift. Not worth blocking.
**50-74%**: Needs attention. Missed requirements, inadequate tests, acceptance issues, or code quality concerns.
**Below 50%**: Significant problems. Major requirements missed, no tests for new behavior, acceptance fails, or risky code.

Nit findings carry zero weight on the overall score — they're stylistic suggestions only and never block a recommendation of Approve.

### Step 9: Report

**Pipeline mode:** Report to the work action orchestrator (which reports to the user).
**Standalone mode:** Report directly to the user.

Format:

```
## Review: REQ-NNN

**Overall: [X]%** | Route [A/B/C] | [commit hash or "uncommitted"]

### Scores

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 95% | All requirements implemented |
| Code Quality | 85% | Clean, follows patterns |
| Test Adequacy | 70% | Missing edge case coverage |
| Scope | 100% | Focused, no drift |
| Risk | None | — |
| Acceptance | Pass | Feature works end-to-end |

### Requirements Checklist

- [x] Dark mode toggle in settings — delivered
- [x] Persists preference in localStorage — delivered
- [ ] Applies to sidebar — not delivered (only main content area)
- [x] Respects OS preference on first visit — delivered

### Findings

**Important:**
- [Specific finding with file:line reference]

**Minor:**
- [Style nit or suggestion]

**Nit:**
- [Optional stylistic suggestion — no score impact]

### Acceptance Testing

**Result: [Pass/Partial/Fail/Untested]**
- [What was tested and outcome]
- [Any issues found during hands-on testing]

### Suggested Additional Testing

- [Manual verification: check dark mode renders correctly across all page layouts]
- [Browser testing: verify localStorage persistence in Safari private mode]
- [Edge case: toggle rapidly between modes to check for flicker/race conditions]

### Follow-up REQs Created
- REQ-025: "Add edge case tests for dark mode toggle" (addendum_to: REQ-003)
```

### Step 9.5: Human Validation (Standalone Mode Only)

In **Standalone mode**, after presenting your review report, you MUST use your environment's ask-user prompt/tool to request human validation.

Prompt the user: *"I have completed my automated review. Please test the deliverable manually. Do you approve? Do you have any feedback, lessons learned, or bugs to report?"*

Based on their response:

- **Lessons Learned / Architectural Feedback:** Append their exact insights directly into the `## Lessons Learned` section of the archived REQ file. (Create the section if it doesn't exist).
- **Bugs / Fixes requested:** Treat these as **Important** findings. Pass them to Step 10 so the system automatically generates new follow-up REQ files (status: pending) containing their exact feedback, linking back via `addendum_to`, so the builder can fix them in the next run.
- **Approved:** Note the approval and proceed.

In **Pipeline mode**, skip this step entirely so you do not block the automated background work loop.

### Step 10: Create Follow-up REQs

For each **Important** finding, create a follow-up REQ:

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
Found during review of [REQ-id]. [1 sentence on what the review found.]

## Requirements
- [Specific fix needed]
```

**When the root cause is ambiguous requirements** — not a code quality issue or missed implementation, but genuine ambiguity in what the user wanted — add an `## Open Questions` section to the follow-up REQ and set its status to `pending-answers`:

```markdown
---
status: pending-answers
---

## Open Questions
- [ ] [What needs clarification before this fix can be implemented]
  Recommended: [best default based on review findings]
  Also: [alternative A], [alternative B]
```

The `pending-answers` status means the work loop won't pick this up until the user reviews it, answers the questions, and flips the status to `pending`. The recommended choices let the user quickly pick an option without deep context-switching. Only add Open Questions when the ambiguity caused the issue — if the fix is clear (e.g., "missed a null check"), use `status: pending` and skip the Open Questions.

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
| Acceptance | [result] |

**Findings:** [count] important, [count] minor
**Acceptance:** [Pass/Partial/Fail/Untested] — [1-line summary]
**Suggested testing:** [count] items
**Follow-ups created:** [REQ-NNN, REQ-NNN] or "None"

*Reviewed by review work action*
```

In standalone mode, this is an exception to the archive immutability rule — review annotations are post-work metadata, not content changes.

### Commit (Standalone mode, git repos only)

In **standalone mode**, after appending the Review section and creating any follow-up REQs, commit the changes. In **pipeline mode**, skip this — the work action's Step 9 handles the commit.

Check for git with `git rev-parse --git-dir 2>/dev/null`. If not a git repo, skip.

```bash
# Stage the modified archived REQ (with appended Review section)
git add do-work/archive/UR-NNN/REQ-NNN-slug.md

# Stage any follow-up REQs created
git add do-work/REQ-NNN-slug.md

git commit -m "$(cat <<'EOF'
[REQ-NNN] review: {score}% (Route {route})

Reviewed: do-work/archive/UR-NNN/REQ-NNN-slug.md
Follow-ups: {REQ-NNN, REQ-NNN} or "None"

EOF
)"
```

**Format:** `[REQ-NNN] review: {score}% (Route {route})` — where `{score}` is the overall review percentage and `{route}` is the original triage route. List the reviewed file path and any follow-up REQs created.

Do not use `git add -A` or `git add .` — stage only the modified archived REQ and any new follow-up REQs. Don't bypass pre-commit hooks.

## Calibrating Review Depth

Match effort to complexity:

| Route | Review depth |
|-------|-------------|
| **A** (Simple) | Quick scan. Requirements check is a glance, code review focuses on correctness, acceptance is a quick smoke test. Skip dimensions that don't apply (e.g., Test Adequacy for a copy change). Suggested testing is usually empty or 1 item. |
| **B** (Medium) | Standard review. Full requirements checklist. Code review checks all dimensions. Acceptance tests the feature end-to-end. Suggested testing covers obvious gaps. |
| **C** (Complex) | Thorough review. Requirements checklist cross-referenced against plan and UR. Code review checks architectural decisions. Acceptance tests multiple paths. Suggested testing is comprehensive — integration, edge cases, performance. |

## What NOT to Do

- Don't re-implement — you're reviewing, not building
- Don't review your own review's follow-up REQs more strictly than the original work — avoid infinite loops of diminishing-return fixes
- Don't block on minor issues — report them but keep moving
- Don't invent requirements — review against what the REQ says, not what you think it should say
- Don't penalize the absence of things the project doesn't have (no test infrastructure = don't fail on test adequacy)
- Don't turn acceptance testing into a full QA cycle — it's a smoke test, not an exhaustive regression suite
- Don't suggest testing for things that are clearly irrelevant to the change
