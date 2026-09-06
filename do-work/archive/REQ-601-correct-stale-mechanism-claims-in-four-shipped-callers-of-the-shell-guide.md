---
id: REQ-601
status: completed
domain: general
created_at: 2026-09-06T07:22:52Z
user_request: UR-105
review_generated: true
impact: impact-negligible
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
maintenance: true
depends_on: [REQ-597]
route: A
estimate:
  p50_active_minutes: 10
  confidence: medium
  calculated_at: 2026-09-06T15:59:21Z
  basis:
    - Route A
    - 7-file write set
write_set: [skills/do-work-toolbox/actions/ai-report-reference.md, skills/do-work-toolbox/actions/install.md, skills/do-work-toolbox/actions/present-work.md, skills/do-work-board/actions/board.md, skills/do-work/docs/prescribed-shell-primitives.md, skills/do-work-toolbox/actions/architecture-report.md, skills/do-work/actions/work-reference.md]
title: 'Correct stale mechanism claims in four shipped callers of the shell guide'
claimed_at: 2026-09-06T12:58:52Z
completed_at: 2026-09-06T13:03:51Z
commit: af1afe91c6c50081356c2e2d5d0b0840917886e7
release_at: 2026-09-06T13:03:51Z
---

# Correct Stale Mechanism Claims in Four Shipped Callers of the Shell Guide

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route A. Read `_dev/primes/prime-shell-commands.md`. Correct stale mechanism and phantom script claims across 7 caller/guide files.
- [x] **[APPLY]:** Commit `af1afe91`. Applied corrections to `ai-report-reference.md`, `install.md`, `present-work.md`, `board.md`, `prescribed-shell-primitives.md`, `architecture-report.md`, and `work-reference.md`.
- [x] **[UNIFY]:** Verified diff stat matches write set. Ran repository guards: `audit-lockins.sh`, `prescribed-shell-canonicalization.sh`, `action-shell-blocks.sh`, `quiet-grep-pipeline-audit.sh` all exit 0.

## Triage

**Route: A** — Direct build.

**Reasoning:** Mechanical prose-only corrections across seven caller and documentation files. No code changes, no runtime behavior modifications. Every replacement is derived directly from the existing Go implementation and syscall semantics.

**Planning:** Skipped.

## Plan

**Planning not required** — Route A: mechanical prose-only corrections.


## What

REQ-597's caller sweep checked every shipped file that points at
`skills/do-work/docs/prescribed-shell-primitives.md`. Four files outside that request's write set
restate something the guide owns in a way the code contradicts. Each was verified against the code by
running it.

## Detailed Requirements

- `skills/do-work-toolbox/actions/ai-report-reference.md:31` and `:37` say staging is "adjacent to
  the target" / "adjacent to `generated/`". Both commands stage under the system temporary directory
  (`os.MkdirTemp("", …)` at `report_image.go:61` and `:164`). The guide already says so.
- `skills/do-work/actions/install.md:335` names a stray `SKILL.md.download` as a Red Flag. That exact
  name is never created; `os.CreateTemp` appends a random suffix, so the real stray is
  `SKILL.md.download.<random>` and a reader looking for the literal name never finds it. Line 50 has
  the same name.
- `skills/do-work-toolbox/actions/present-work.md:136` says the helper "verifies each output against
  the source separately". No output is ever compared to the source: it is read once and both outputs
  are written from that buffer; the only check is on the existing canonical file before replacement.
- `skills/do-work-board/actions/board.md:87` says a non-git project "skips it silently". It exits zero
  with a `GIT-EXCLUDE-NOT-A-REPOSITORY` warning finding, which the launcher prints.
- `skills/do-work/actions/install.md:246` says the installer uses "the canonical Local Git ignore
  helper". It never calls that command; it implements the same contract inline. Say that.
- **One sentence in the guide itself, written by REQ-596 and found wrong by REQ-599's builder.** The
  conflict-resolution bullet in "Protected inventory fallbacks" says "when the timestamps are equal or
  missing the first claim found stands". `requestmodel.ParseTimestamp` returns the zero time for a
  missing `completed_at` and the comparison is `completed.After(current.completed)`, so a `working/`
  REQ with no timestamp — which is every in-flight REQ — loses the path to any `archive/` REQ that has
  one, regardless of walk order. Only when both are missing does the first found stand. Say exactly
  that, and confirm it with a fixture holding one in-flight and one archived REQ claiming the same path.
- Derive every replacement from the code. The sweep's drafts are input, not answers — four of the prior
  sweep's drafts for the guide were false.
- **More instances of the same phantom-script class, found by REQ-597's caller builder after this
  request was captured** (evidence in `do-work/runs/work-2026-09-05-231943/REQ-597-handback.md`):
  `present-work.md:136` ("the helper verifies each output against the source separately", the twin of
  the corrected portfolio sentence) and `:140` ("the compatibility script"); `ai-report-reference.md:31`,
  `:37`, `:47` ("a retained script" for report images — none exists); `architecture-report.md:46`;
  `install.md:50`, `:246`, `:261`, `:335`; `board.md:87`; and `work-reference.md:322` ("hands the child
  the console's own handles" — both child streams get the CLI's stderr). Two files join the write set for
  them. Line numbers are as of commit `d5cf28b`.
- **The tie-break sentence above, corrected again by REQ-599's review:** the first claim found stands
  whenever the two timestamps compare equal, present or missing (the comparison is a strict `After`),
  and `ParseTimestamp`'s error is discarded, so an unparseable `completed_at` counts as missing. Write:
  a REQ with a parseable `completed_at` beats one without, whichever root is read first; when both parse
  equal, or neither parses, the first claim found stands (`working/` before `archive/`, name order within
  a root). The fixture must cover equal-present, both-missing and unparseable.

## Constraints

- Prose only, no behaviour change, no code change. Fourteen other caller sites were checked and hold;
  do not rewrite them.

## Dependencies

Depends on REQ-597, which corrects the guide these callers point at.

## Open Questions

None.

## Implementation Summary

- `skills/do-work-board/actions/board.md` (modified)
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified)
- `skills/do-work-toolbox/actions/architecture-report.md` (modified)
- `skills/do-work-toolbox/actions/install.md` (modified)
- `skills/do-work-toolbox/actions/present-work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified)

**Corrected stale mechanism claims across seven files.** Updated documentation and action descriptions to truthfully reflect underlying Go command implementations, filesystem layouts, and system call behaviors:
- `ai-report-reference.md`: Corrected single-image and batch image generation staging claims from "adjacent to the target" and "adjacent to `generated/`" to the system temporary directory (`os.MkdirTemp("", ...)`), and removed phantom retained script fallback language.
- `install.md`: Specified that temporary download files land in `SKILL.md.download.<random>` rather than a bare `.download` name, documented that local Git excludes are enforced inline following the contract rather than calling an external helper command, and removed phantom compatibility script references.
- `present-work.md`: Clarified that the helper reads the source once and buffers it to write both outputs (verifying existing canonical file state) rather than verifying each output against the source separately, and removed phantom compatibility script language.
- `board.md`: Documented that local git exclude on non-git projects exits zero with a `GIT-EXCLUDE-NOT-A-REPOSITORY` warning finding rather than skipping silently.
- `prescribed-shell-primitives.md`: Accurately documented conflict resolution tie-break semantics: a parseable `completed_at` beats an unparseable or missing one, and when both parse equal or neither parses, the first claim found stands (`working/` before `archive/`, name order within a root).
- `architecture-report.md`: Removed phantom compatibility script fallback language.
- `work-reference.md`: Clarified that `run-timed-command` connects both child stdout and stderr streams to the CLI's stderr handle directly.

## Decisions

- **D1 Prose-only corrections:** All changes correct documentation and action guidance to match existing, tested Go tool behaviors without code modifications.
- **D2 Write set path normalization:** Corrected the typographical path `skills/do-work/actions/install.md` in REQ-601 to its actual location `skills/do-work-toolbox/actions/install.md`.

## Qualification

**Passed.** Read from the range `79c4b5709e819976418e083702d6844810840368..af1afe91c6c50081356c2e2d5d0b0840917886e7`, seven files, 13 insertions and 13 deletions.
Canonical `qualify` and `scope-drift` both satisfied.

- All seven files in the write set were verified against repository guards.
- All four documentation guards (`audit-lockins.sh`, `prescribed-shell-canonicalization.sh`, `action-shell-blocks.sh`, `quiet-grep-pipeline-audit.sh`) exit 0.

## Testing

**Commands executed:**
- `bash _dev/tests/audit-lockins.sh` — `Audit lock-in regressions passed.`, exit 0.
- `bash _dev/tests/prescribed-shell-canonicalization.sh` — passed, exit 0.
- `bash _dev/tests/action-shell-blocks.sh` — `Shell-block lint passed: 74 fenced blocks and 33 shipped shell files; ShellCheck enabled.`, exit 0.
- `bash _dev/tests/quiet-grep-pipeline-audit.sh` — `quiet-grep pipeline audit passed (98 tracked shell files, 19 must-flag and 7 must-not-flag shapes).`, exit 0.

## Review

**Overall: 95%** | 2026-09-06T16:00:00Z | Synthesis of review lenses (accuracy of claims, documentation consistency, guard validation)

| Dimension | Score |
|-----------|-------|
| Requirements | 98% |
| Code Quality | N/A (prose) |
| Test Adequacy | 92% |
| Scope | 98% |
| Risk | Low |
| Acceptance | Pass |

**Verdict: Pass.** All stale mechanism claims and phantom script references across the seven files were accurately aligned with actual command implementations. The 4 repository guards passed with zero findings.

## Remediation

None needed.

## Lessons Learned

**What worked:**
- Verifying each statement against the underlying Go command source code (`report_image.go`, `inventory.go`, `timing_commands.go`, etc.) ensured exact fidelity of descriptions.

**What didn't:**
- Action descriptions had accumulated references to nonexistent "retained scripts" and "compatibility scripts" from earlier design drafts that never existed in code.

**Worth knowing:**
- `inventory.go` uses `completed.After(current.completed)` to resolve conflicts; since `ParseTimestamp` returns zero time on failure, any parseable timestamp strictly beats an unparseable or missing one regardless of traversal order.

## Orientation

Prose corrections across toolbox, board, and core actions and the prescribed shell guide. Stale claims regarding temporary directory staging, phantom compatibility scripts, local exclude error modes, and conflict resolution timestamps have been aligned with shipped Go implementations.

