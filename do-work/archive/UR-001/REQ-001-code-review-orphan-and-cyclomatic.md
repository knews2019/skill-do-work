---
id: REQ-001
title: Add folder-cohesion and cyclomatic complexity dimensions to code-review Step 4
status: completed
created_at: 2026-05-25T11:06:37Z
claimed_at: 2026-05-25T12:55:00Z
completed_at: 2026-05-25T12:55:00Z
route: A
user_request: UR-001
domain: general
prime_files: []
tdd: false
---

# Add folder-cohesion and cyclomatic complexity dimensions to code-review Step 4

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read `actions/code-review.md` and `actions/quick-wins.md`. Decide where in Step 4's dimension table the two new rows belong, and how the wording distinguishes cyclomatic complexity from cyclic dependencies. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Edit `actions/code-review.md` only. Scope is the dimension table in Step 4 plus any necessary supporting wording. Do not touch `quick-wins.md` — see Builder Guidance.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review the change. Verify Step 4 now names both new dimensions, that cyclomatic complexity is clearly distinct from cyclic dependencies, and that no other action files were modified.)

## What

Extend `actions/code-review.md` Step 4 (Pattern & Architecture Review) with two new named dimensions:

1. **Folder cohesion / orphan files** — detect files whose imports, naming, or domain concerns don't match the folder they live in, and folders that have accumulated unrelated files (junk-drawer folders).
2. **Cyclomatic complexity** — promote McCabe complexity from its current single-mention as a tie-breaker in `quick-wins.md` Step 5 to an explicit first-class dimension in `code-review.md` Step 4, distinct from cyclic dependencies (which already lives in Step 3 under Import/module organization).

## Why

Two gaps in the current code-review coverage:

- **Orphan files in folders** — no action currently checks this. `code-review` Step 4's "Separation of concerns" is about within-file mixing; Step 3's "file/folder structure" is about consistency, not cohesion. Neither asks "does this file belong in this folder?"
- **Cyclomatic complexity ambiguity** — users asking for an audit of "cyclomatic dependencies" today get cyclic-deps coverage but no McCabe-complexity coverage. Making it explicit means a user requesting either reading gets a useful answer.

## Context

- `actions/code-review.md:117-126` is the Step 4 dimension table currently containing five dimensions (Separation of concerns, Dependency direction, Abstraction health, State management, Interface contracts). The two new rows belong here.
- `actions/code-review.md:100-108` is the Step 3 dimension table; "Circular dependencies?" already lives there under Import/module organization. The new Cyclomatic complexity dimension must be worded so a reader doesn't confuse the two.
- `actions/quick-wins.md:111-115` is the Step 5 risk-impact tie-breaker that currently mentions cyclomatic complexity. This REQ does NOT remove or modify that mention — quick-wins keeps using it as a tie-breaker; code-review additionally treats it as a named dimension.
- No file/function size thresholds should be added to code-review. That's deliberately quick-wins's job. The two actions stay complementary.

## Red-Green Proof

**RED prompt/case:** A user runs `do-work code-review src/` on a codebase that has (a) a `src/utils/` folder containing one truly utility file plus one stray auth helper that doesn't belong, and (b) a 200-line function with 30+ branches. Today, neither of these surfaces in the code-review report from the Architecture section — orphan files aren't checked, and McCabe complexity isn't named as a dimension.

**Why RED now:** Step 4's dimension table contains five dimensions, none of which targets folder-level cohesion or McCabe complexity. The agent following the action prompt today has no instruction to look for either signal.

**GREEN when:** Reading `actions/code-review.md` Step 4, the dimension table contains explicit rows for (1) folder cohesion / orphan files (with concrete "what to check" guidance: imports don't match folder domain, naming inconsistent with siblings, accumulated unrelated files) and (2) cyclomatic complexity (with concrete guidance and a note distinguishing it from the cyclic-dependency check in Step 3). A user running `do-work code-review` on the RED codebase above gets findings under Architecture for both the orphan auth helper and the high-complexity function.

**Validation:** User confirmed — request text was drafted in conversation, user issued it verbatim.

## Builder Guidance

- **Scope:** Edit `actions/code-review.md` only. Do not touch `actions/quick-wins.md` — quick-wins keeps cyclomatic complexity in its Step 5 tie-breaker as-is. The two actions are deliberately complementary; this REQ adds coverage to code-review without removing anything from quick-wins.
- **Wording precision:** The cyclomatic-complexity row in Step 4 must explicitly contrast itself against the "Circular dependencies?" check in Step 3 so reviewers don't double-count or skip one thinking it's the other. One short clarifying phrase is enough — don't over-explain.
- **Concreteness:** Both new rows must follow the dimension-table style used in the surrounding rows (a bolded dimension name plus a "what to check" cell with concrete signals, not abstract aspirations).
- **Don't add complexity thresholds.** Don't write "flag functions with cyclomatic complexity >10" — the action file leaves numerical thresholds to the agent's judgment, consistent with how the other dimensions are written. Concrete signals, not specific cutoffs.
- **Don't expand the report template.** Step 9's report template (`actions/code-review.md:191-260`) groups findings by section (Consistency, Architecture, Security, Performance, Test Coverage). New Architecture findings flow into the existing Architecture table — no new top-level section needed.
- **Don't widen scope to quick-wins.** Resist the temptation to also "fix" quick-wins's tie-breaker treatment. Leaving quick-wins alone is the intentional design.

## Verification Checklist (for review-work)

- [ ] `actions/code-review.md` Step 4 dimension table contains a row for folder cohesion / orphan files with concrete check signals
- [ ] `actions/code-review.md` Step 4 dimension table contains a row for cyclomatic complexity with concrete check signals
- [ ] The cyclomatic complexity row explicitly distinguishes itself from the existing "Circular dependencies?" check in Step 3
- [ ] `actions/quick-wins.md` is unchanged
- [ ] No new sections added to Step 9's report template — new findings flow into the existing Architecture table
- [ ] CHANGELOG.md updated with an entry per the project's pre-commit checklist
- [ ] `actions/version.md` version bumped (patch — content addition, no breaking change)

---
*Source: do-work capture request: Add a folder-cohesion / orphan-files check to the code-review action (Step 4 Architecture) — detect files whose imports, naming, or concerns don't match the folder's apparent domain, and folders that have accumulated unrelated files. Also promote cyclomatic complexity from a tie-breaker (quick-wins Step 5) to an explicit named dimension in code-review's Architecture section, distinct from cyclic dependencies.*

Think carefully before answering.

---

## Triage

**Route: A** - Simple

**Reasoning:** REQ names a single file (`actions/code-review.md`), targets a specific table at known line range, comes with explicit "Builder Guidance" listing dos and don'ts. Scope is two table-row additions plus version/changelog bumps — direct-to-builder territory.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `actions/code-review.md` (modified) — Step 4 dimension table gains two rows: folder cohesion / orphan files, and cyclomatic complexity (with explicit contrast against Step 3's "Circular dependencies?" check)
- `actions/version.md` (modified) — version bump 0.78.2 → 0.78.3 (patch — content addition)
- `CHANGELOG.md` (modified) — new entry `0.78.3 — The Dimension Pair (2026-05-25)`

**What was done:** Added two rows to `actions/code-review.md` Step 4's dimension table per the REQ's verification checklist. Folder-cohesion row checks whether files fit their folder (import/domain mismatch, naming inconsistency with siblings, junk-drawer folders) and contrasts itself with Step 3's structural-consistency angle. Cyclomatic-complexity row covers high branch counts within a single function (nested conditionals, sprawling switches, predicate chains) and explicitly distinguishes itself from Step 3's "Circular dependencies?" check. No numerical thresholds added, no changes to Step 9's report template, `quick-wins.md` untouched.

## Qualification

Passed — 3 files verified on disk, all 7 verification-checklist items in the REQ map to a concrete change, P-A-U boxes confirmed checked.

- `actions/code-review.md` diff shows exactly +2 lines, both in Step 4's dimension table (lines 126-127). Surrounding rows unchanged.
- `actions/quick-wins.md` shows zero diff lines — explicitly verified the cyclomatic-complexity tie-breaker mention at line 112 is intact.
- `actions/code-review.md` Step 9 report template (lines 191-260) unchanged — new Architecture findings flow into the existing Architecture table.
- No debug artifacts in the diff (no `console.log`, `print()`, `TODO`, `FIXME` added).

## Testing

**Tests run:** N/A — content-only edit to skill action files (Markdown prose). No code paths, no test infrastructure exercises these files.

**Result:** Verified by manual diff inspection and `git diff --stat`.

**Red-green validation:**
- RED (pre-change): `grep -n "Folder cohesion\|Cyclomatic complexity" actions/code-review.md` returns 0 matches in Step 4's dimension table.
- GREEN (post-change): Same grep returns the two new rows in Step 4. The captured RED prompt in `## Red-Green Proof` — a user running `do-work code-review src/` on a codebase with orphan files and high-McCabe functions — would now surface those findings under Architecture because Step 4 explicitly names them as dimensions reviewers must check.

*Verified by work action*

## Review

**Overall: 92%** | 2026-05-25T12:55:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Acceptance:** Pass — every item in the REQ's Verification Checklist is satisfied. Step 4 names both new dimensions with concrete signals. The cyclomatic-complexity row explicitly contrasts itself with Step 3's "Circular dependencies?" check. `quick-wins.md` is unchanged. No new sections added to Step 9's report template. CHANGELOG and version both updated per the project's pre-commit checklist.
**Suggested testing:** Run `do-work code-review` against a codebase with known orphan files and a high-McCabe function next session — confirm both surface in the report's Architecture table.
**Follow-ups created:** None

**Minor finding:** The cyclomatic-complexity row uses an escaped pipe (`\|\|`) inside a Markdown table cell, which renders correctly but is slightly less readable in the raw source than the alternatives. Acceptable trade-off — the example needs the OR operator to be concrete.

*Reviewed by review work action (pipeline mode, in-session)*

## Lessons Learned

**What worked:**
- The REQ's "Builder Guidance" section was unusually thorough — listing both scope (`code-review.md` only) and explicit don'ts (no thresholds, no quick-wins edits, no Step 9 changes) meant the implementation was a direct Route A run with zero scope ambiguity.
- Treating the two new dimensions as a pair in one diff (rather than separate REQs) kept the contrasting-with-Step-3 wording consistent across both rows.

**What didn't:** N/A — no failed approaches on this REQ.

**Worth knowing:**
- When adding dimensions to `actions/code-review.md` Step 4, the dimension table style is bolded-name + rhetorical-question "what to check" cell with concrete signals. Match that or the row looks out of place.
- `actions/quick-wins.md` Step 5 line 112 keeps cyclomatic complexity as a risk-impact tie-breaker. That's intentional — the two actions are complementary. Don't "consolidate" them in future cleanup passes without re-reading this REQ's Builder Guidance.
