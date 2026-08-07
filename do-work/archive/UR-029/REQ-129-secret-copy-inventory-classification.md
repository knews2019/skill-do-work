---
id: REQ-129
title: Secret-derived copies remain quarantined
status: completed
claimed_at: 2026-08-07T10:52:04Z
completed_at: 2026-08-07T10:56:44Z
commit: 4ea68ae
route: B
created_at: 2026-08-07T08:45:11Z
user_request: UR-029
addendum_to: REQ-128
domain: security
prime_files: []
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-130, REQ-131, REQ-132]
batch: audited-safety-fixes
effort_estimate: normal
write_set: [tools/checks/uncommitted-inventory.sh, actions/commit.md, actions/inspect.md, _dev/tests/contract-regressions.sh, actions/version.md, CHANGELOG.md]
---

# Addendum: Secret-Derived Copies Remain Quarantined

## What

Complete the secret-safe inventory work by preserving Git copy records and applying secret-shaped basename classification to both their source and destination. Keep REQ-128's stronger rename quarantine and fail-closed re-inventory behavior intact.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

The deep audit proved that the current parser consumes the second NUL-delimited field and correctly emits `XD` for a secret rename source plus `X` for its destination. An ordinary rename remains `M`. The remaining defect is command-level: when Git is configured for copy detection, the explicit rename-only option downgrades an emitted `C` record to `A`, so copying a modified `.env.local` to `config.txt` leaves the destination readable.

## Prior Implementation

REQ-128 hardened REQ-121's inventory workflow in commit `7bb03d2`. It forced rename detection, retained a run-level once-`X`-always-`X` quarantine, converted ambiguous additions beside `XD` to `X`, and safely handled already-staged secret deletions. Its files and tests are the baseline; do not remove or weaken those protections.

## Detailed Requirements

- Make the porcelain invocation explicitly versioned, NUL-delimited, exhaustive for untracked files, and copy-aware while overriding repository configuration that disables rename detection.
- Preserve and parse both path fields for every emitted rename or copy record.
- Classify a copy destination as `X` when either its source or destination basename is secret-shaped.
- Preserve secret rename output as `XD` for the source and `X` for the destination; do not reduce it to the consumer report's older single-row expectation.
- Preserve an ordinary rename as `M` and an unrelated ordinary addition as `A` when no ambiguous `XD` state exists.
- Mirror the copy-aware command and both-path classification in the commit and inspect manual fallbacks.
- Preserve Bash 3.2 compatibility and the public `M`, `A`, `D`, `X`, and `XD` interface.
- Do not introduce content-based copy detection or read secret contents.

## Constraints

- Preserve all completed behavior from REQ-128, including rename-disabled and reset-and-reinventory protection.
- Add tests in the existing shell regression framework; do not introduce consumer-specific Jest coverage.
- Run relevant regressions, `bash -n`, ShellCheck, and Bash 3.2 checks. Run the repository's broader validation required by its contribution rules.
- Do not stage, commit, or push without separate user authorization.

## Builder Guidance

Firm on behavior and safety. The audit showed that forcing copy-aware status configuration preserves `C` records even when repository configuration says renames are disabled; the builder may choose the cleanest portable spelling that proves both conditions without weakening REQ-128.

## Red-Green Proof

**RED prompt/case:** In a temporary Git repository, commit `.env.local`, modify it, copy the modified file to `config.txt`, stage both paths, enable copy detection, and run the inventory. Git can emit a `C` record, but the current inventory command downgrades `config.txt` to `A`. Separately, stage a secret rename and an ordinary rename.
**Why RED now:** The parser has a defensive `C` branch, but the status invocation suppresses the copy record before the parser can classify its source.
**GREEN when:** The copy destination is `X`; the secret rename remains `XD .env.local` plus `X config.txt`; the ordinary rename remains `M`; both action fallbacks prescribe the same behavior; and the regression fails if copy detection is downgraded again.
**Validation:** User confirmed by requesting capture of the audited findings after reviewing the reproduced evidence.

## Dependencies

No execution dependency on another queued REQ. This addendum corrects completed REQ-128.

## Full Context

See `do-work/user-requests/UR-029/input.md` for the complete capture context.

---
*Source: UR-029 - "run do-work capture-request on these issues"*

---

## Triage

**Route: B** - Medium

**Reasoning:** The root cause is localized, but the fix must synchronize one executable checker, two prose fallbacks, and behavior-level regression coverage.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

The checker already parses the second NUL field for both rename and copy records, but its explicit rename-only status option prevents Git from emitting copy records. Both action fallbacks repeat that option. Existing behavior probes cover secret renames, disabled rename configuration, and ambiguous additions, so the new fixture can isolate copy classification while retaining ordinary rename behavior.

## Scope

**Files I will touch:**
- `tools/checks/uncommitted-inventory.sh` (modify) - request copy-aware porcelain records while preserving defensive parsing
- `actions/commit.md` (modify) - synchronize the manual fallback command and copy contract
- `actions/inspect.md` (modify) - synchronize the manual fallback command and copy contract
- `_dev/tests/contract-regressions.sh` (modify) - add copy and ordinary-rename behavior probes
- `actions/version.md` (modify) - bump the integration version
- `CHANGELOG.md` (modify) - record the safety fix

**Files I will NOT touch:** Other inventory callers or archived antecedents.

**Acceptance criteria (restated from REQ):**
- [ ] A destination copied from a secret-shaped tracked source is tagged X instead of A.
- [ ] A secret rename remains XD for its source and X for its destination.
- [ ] An ordinary rename remains M.
- [ ] Repository rename configuration cannot disable copy-aware detection.
- [ ] Both manual fallbacks prescribe the same copy-aware command.

## Implementation Summary

**Files changed:**
- `tools/checks/uncommitted-inventory.sh` (modified)
- `actions/commit.md` (modified)
- `actions/inspect.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)
- `actions/version.md` (modified)
- `CHANGELOG.md` (modified)

**What was done:** Forced copy-aware porcelain detection independently of repository configuration, synchronized both manual fallbacks, and added behavior probes for secret-derived copies and ordinary renames.

## Qualification

Passed - 6 files verified, 5 acceptance requirements traced, P-A-U confirmed.

## Testing

**Tests run:**
- /bin/bash _dev/tests/contract-regressions.sh - passed

**Red-green validation:**
- RED: The new copy fixture reported copied-config.txt as A, and both fallback contract checks failed on their rename-only commands.
- GREEN: The full contract suite passed after forcing copy-aware porcelain detection and synchronizing both fallbacks.

**Existing tests updated:** _dev/tests/contract-regressions.sh retains the secret rename XD plus X assertions and adds secret-copy X plus ordinary-rename M assertions.

## Root Cause

The parser understood copy records, but its explicit --renames option prevented Git from emitting them and overrode any copy-aware repository configuration. The manual fallbacks repeated the same command-level blind spot.

## Review

**Overall: 100%** | 2026-08-07

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition - this is the durable audit record the gate mandates):**
- None

**Minor findings:** 0 (report only)
**Acceptance:** Pass - copy, secret rename, ordinary rename, and fallback behavior all pass the Bash 3.2 regression suite.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** A dedicated staged-copy fixture exposed the command-level detection mode without weakening existing quarantine cases.

**What didn't:** A tiny unstaged source was not an eligible Git copy candidate, so the initial fixture emitted a truthful addition rather than the intended copy record.

**Worth knowing:** Copy detection needs an eligible modified source with enough similarity; forcing status.renames=copies is independent of repository configuration and preserves ordinary rename records.
