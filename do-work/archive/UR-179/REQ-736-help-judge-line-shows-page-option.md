---
id: REQ-736
title: '[impact-negligible] Show the --page option on the ai-report judge line in help'
status: completed
created_at: 2026-11-15T07:25:23Z
user_request: UR-179
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-releases.md]
tdd: false
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
related: [REQ-735]
batch: run-012821-follow-ups
write_set: [skills/do-work-toolbox/actions/help.md]
claimed_at: 2026-11-15T07:30:05Z
route: A
estimate:
  p50_active_minutes: 10
  confidence: high
  calculated_at: 2026-11-15T07:31:33Z
  basis:
    - Route A
    - 1-file write set
    - 2 acceptance criteria
dispatch_at: 2026-11-15T07:32:53Z
builder_handback_at: 2026-11-15T07:41:03Z
integration_at: 2026-11-15T07:41:12Z
review_at: 2026-11-15T07:47:37Z
commit: e8a4c69c3e190ff6c9b42321534e62dae200a3bd
heavy_verified_at: 2026-11-15T07:51:35Z
heavy_verified_revision: e8a4c69c3e190ff6c9b42321534e62dae200a3bd
kb_status: pending
completed_at: 2026-11-15T07:51:53Z
release_at: 2026-11-15T07:51:53Z
---
# Show the --page Option on the ai-report Judge Line in Help
## What
`skills/do-work-toolbox/actions/help.md` line 18 lists `ai-report judge <dir>` without the `--page` option that REQ-728 (ai-report-judge renders one named page, released 0.319.0) added. Change that one line so it shows the option.
## Why
The maintainer's words: "capture a follow-up REQ for this one line". The run-end restatement sweep of REQ-734 (supports-and-limits table page) found the gap (`do-work/runs/work-2026-11-15-012821/REQ-734-integration-report.md:11`, decision D-20). The line leaves out a form but says nothing false.
## Context
- Current line 18 (verified at capture, commit `db458e9e`): `  ai-report judge <dir>          Render-check a report bundle (widths, themes, links)`.
- The shipped form to match: `judge <bundle-dir> [--page <file.html>]` in `skills/do-work-toolbox/actions/ai-report.md:36`, and `do-work-toolbox ai-report judge <dir> [--page <file.html>]  Render-check an existing report bundle, or one named page of it` in `skills/do-work-toolbox/docs/ai-report-guide.md:80`.
- The line sits in a column-aligned listing (lines 15-20). The builder keeps the listing readable if the longer usage overflows the column.
- Fold-first scan: no queued REQ or sweep shares this root cause (the queue holds only REQ-727, a relation-map preset). The fix is prose-only, but the maintainer explicitly promoted this report-only finding to a REQ, so it does not go to `do-work/prose-backlog.md`.
- Independent of REQ-735 (scope of the moved standing rules). No `depends_on`.
- The change touches a shipped file under `skills/`, so it is a release. Follow `_dev/primes/prime-releases.md`.
## Red-Green Proof
**RED prompt/case:** `sed -n 18p skills/do-work-toolbox/actions/help.md` prints `ai-report judge <dir>` with no `--page`.
**Why RED now:** REQ-728 added `--page` to `ai-report.md` and the guide, but not to the help listing.
**GREEN when:** The judge line in `help.md` shows `[--page <file.html>]`, in the same form as `ai-report.md:36`, and no other line of `help.md` changes.
**Validation:** Inferred during capture
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 8761 tokens, over the 2000-token budget and `slugged: partial`, so it cannot be narrowed. Matching reason: the request changes an action file.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** replace help.md line 18 usage with `ai-report judge <dir> [--page <file.html>]`, keep the description after two spaces (D-01). (from the builder hand-back)
- [x] **[APPLY]:** one exact-string replacement, one commit by exact path (`ab48d0a0`). (from the builder hand-back)
- [x] **[UNIFY]:** `git diff --stat HEAD~1` shows only help.md, 1+/1-; checked the diff holds no other line and no debug text; probe and reference contract pass. (from the builder hand-back)
*Source: item 2 of `do-work/user-requests/UR-179/input.md` (help.md fix after REQ-728)*

---

## Triage

**Route: A** - Simple

**Reasoning:** One prose line in one shipped file (`skills/do-work-toolbox/actions/help.md:18`) gains the `[--page <file.html>]` form that `ai-report.md:36` already documents. No code, no other file, no design choice beyond column alignment.

**Planning:** Not required

**Decisions:**
- **D-01**: The new usage `ai-report judge <dir> [--page <file.html>]` (41 characters) overflows the 30-character usage column. Builder keeps the description on the same line after two spaces, the same shape line 20 (`ai-report --kind proposal|root-cause <target>  Decision brief…`) already uses, and shortens nothing else. Reasoning: matches an existing overflow precedent in the same listing; no other line changes. Value: one-line diff. Risk: none; prose only, trivially reversible.

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions

(from the builder hand-back)

- D-02 (DECIDE & STATE): kept the description unchanged ("Render-check a report bundle (widths, themes, links)"). Adding "or one named page" would make the overflowed line longer, and the usage already shows the option.

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/help.md` (modified, line 18 only)

**What was done:** The `ai-report judge` line in the toolbox help listing now reads `ai-report judge <dir> [--page <file.html>]`, the same form `ai-report.md:36` documents. The longer usage overflows the 30-character column, so the description follows after two spaces, as the `ai-report --kind` line already does (D-01). The description text is unchanged (D-02).

## Discovered Tasks

(from the builder hand-back)

- `_dev/tests/staged-skills-contract.sh` is named as a builder verify step in the brief, but it refuses a standalone run (exit 2, heavy-only). Briefs should drop it or name the heavy lane for the integrator. → report only

## Qualification

**Gate records:** `advance --diff-range 15402396..e8a4c69c` returned phase `qualify`, gate record `qualify` satisfied, no findings.

**Requirement trace against `git diff 15402396..e8a4c69c` (1 file, +1/-1):**
- The judge line now shows `[--page <file.html>]`: `skills/do-work-toolbox/actions/help.md:18` reads `ai-report judge <dir> [--page <file.html>]  Render-check a report bundle (widths, themes, links)`.
- Same form as the source: `skills/do-work-toolbox/actions/ai-report.md:36` documents `judge <bundle-dir> [--page <file.html>]`; the option text matches byte for byte. `<dir>` stays as the listing's existing placeholder, as in `ai-report-guide.md:80`.
- No other line of `help.md` changes: the diff is one hunk, one line removed and one added.
- Column overflow handled as D-01 decided: two spaces before the description, the shape line 20 already uses.

**Scope:** Route A; write_set = `skills/do-work-toolbox/actions/help.md`, and the diff touches only that file. No debug artifacts.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` at merge `e8a4c69c3e190ff6c9b42321534e62dae200a3bd` (detached drain-head checkout, `DO_WORK_FAST_STAGE_REUSE=off`, 1-minute load 4.59 at start); GREEN probe `do-work/runs/work-2026-11-15-012821/REQ-736-probe.sh` through `advance --probe-file`.
**Result:** ✓ Gate exit 0, 168 s, "Maintainer verification passed." Probe exit 0; gate records `run-blocked-check` and `green-gate` satisfied.

**Red-green validation:** (from the builder hand-back)
- `REQ-736-probe.sh`: ✗ at base `547d4cc5` ("FAIL: help.md judge line does not show [--page <file.html>]", exit 1) → ✓ at tip `ab48d0a0` ("PASS: REQ-736 help.md judge line shows --page", exit 0)

**New tests added:**
- None (one prose line; the run probe is run evidence, not a shipped test).

**Heavy verification plan:** `plan-heavy-verification` over 154023962e969085af278b8d4e4e840c86b22f05..e8a4c69c3e190ff6c9b42321534e62dae200a3bd selected one lane:
- `staged-skills` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — reason: `skills/do-work-toolbox/actions/help.md matched subtree skills`. The brief also required this lane, because the builder could not run `staged-skills-contract.sh` standalone.

## Review

**Overall: 98%** | 2026-11-15T07:47:37Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 98% |
| Test Adequacy | N/A (one prose line; run probe red at 547d4cc5, green at ab48d0a0) |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** None
**Acceptance:** Pass — `sed -n 18p skills/do-work-toolbox/actions/help.md` shows `[--page <file.html>]` in the ai-report.md:36 form.
**Restatement sweep:** nothing redefined
**Suggested testing:** 0 items
**Follow-ups created:** None (0 findings report only)
**Anti-bloat:** 0 helpers, options, files or decorative tests beyond the REQ (reviewer count).
**Full report:** `do-work/runs/work-2026-11-15-012821/REQ-736-review.md`

*Reviewed by review-work action*

## Heavy Verification Plan

- Base: `154023962e969085af278b8d4e4e840c86b22f05`
- Target: `e8a4c69c3e190ff6c9b42321534e62dae200a3bd`
- `staged-skills` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — reasons: `skills/do-work-toolbox/actions/help.md matched subtree skills`

## Heavy Verification Result

- Target: `e8a4c69c3e190ff6c9b42321534e62dae200a3bd`; execution revision `e8a4c69c3e190ff6c9b42321534e62dae200a3bd` (detached drain-head checkout, `QUEUE_KANBAN_BROWSER` set).
- `staged-skills` — executed, exit 0, 38 s.

## Lessons Learned

**What worked:** None beyond the hand-back: a one-line prose fix with a byte-level probe against the shipped source form (`ai-report.md:36`) left nothing to judge at integration.
**What didn't:** The builder brief named `_dev/tests/staged-skills-contract.sh` as a builder step, but it runs only as a heavy lane (exit 2 standalone); the integrator's `staged-skills` lane covered it.
**Worth knowing:** A listing usage longer than the help column follows the line-20 precedent: two spaces, then the description.

## Orientation

`skills/do-work-toolbox/actions/help.md` is the toolbox's column-aligned command menu; its usages restate forms owned by each action file, so a new option in an action file needs the matching help line.

## Timing

Observed 2026-11-15T07:31:00Z to 2026-11-15T07:51:35Z: 20m 35s total, 14m 29s attributed across 6 events, 6m 06s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 8m 11s | 1 |
| verification-gate | 3m 44s | 2 |
| planning | 1m 36s | 1 |
| review | 49s | 1 |
| handback-merge | 9s | 1 |

Slowest stage: builder-work / builder worktree build, 8m 11s, outcome success.
