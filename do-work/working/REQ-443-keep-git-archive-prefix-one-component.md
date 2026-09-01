---
id: REQ-443
title: '[impact-critical] Keep Git fallback archive prefixes to one component'
status: claimed
route: A
created_at: 2026-08-31T14:19:37Z
user_request: UR-083
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-01T22:16:10Z
  basis:
    - regression coverage for an already-centralized constant
related: [REQ-437, REQ-438, REQ-439, REQ-440, REQ-441, REQ-442, REQ-444]
batch: accepted-feedback-regressions
claimed_at: 2026-09-01T22:14:46Z
---

# Keep Git Fallback Archive Prefixes to One Component

## What

Use a constant single-component prefix for Git-fallback archives regardless of the selected branch name. Branch text must select the cloned ref only; it must not alter the extraction depth expected by install and update.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add one production-style regression that forces HTTP fallback, clones an exact slash-containing branch, extracts with `--strip-components=1`, and asserts that branch's `VERSION` at the extraction root. The constant-prefix implementation is already present from `2c82ef12`, so prove the test detects the historical branch-derived prefix with a controlled mutation and retain production code unchanged.
- [x] **[APPLY]:** Added the slash-branch Git fallback regression; retained the constant `upstream/` production prefix already introduced by `2c82ef12`.
- [x] **[UNIFY]:** Reviewed `archive_fetch_test.go` and the request record in scoped diff context; verified exact branch selection, failed-HTTP fallback, one-component production extraction, root `VERSION`, and no debug artifacts. The production file has no retained diff.

## Finding Provenance

- **Verbatim claim / severity:** `[P2] Keep Git archive prefixes to one path component.`
- **Evidence:** `archive_fetch.go:108-110` embeds the branch name in `repo-<branch>/`, while install and update extract with `--strip-components=1`.
- **Origin / earned by:** Shell commit `0e8cf0d9` introduced the prefix shape, `f27f564d` preserved it, and exact branch selection in REQ-424 (Clone the Branch Named by the Tarball URL) made slash-containing refs load-bearing. A `release/2.x` replay extracted `VERSION` below `2.x/`, so root manifest validation failed.
- **Surface-cost:** N/A. A constant prefix deletes the incorrect coupling and adds no defensive apparatus.

## Detailed Requirements

- Keep derived branch names as exact Git clone selectors, including names containing `/`.
- Generate every Git-fallback archive beneath one constant path component.
- Preserve the production extraction contract that strips exactly one component.
- Cover a slash-containing local fixture branch through fetch, archive, production-style extraction, and root manifest-file assertion.

## Constraints

- Do not sanitize branch text into a second naming scheme; remove it from the prefix entirely.
- Preserve missing-branch failure and exact requested-branch selection from REQ-424.

## Red-Green Proof

**RED prompt/case:** Force HTTP failure, fetch a local Git branch named `release/2.x`, extract with production's `--strip-components=1`, and look for `VERSION` at the extraction root.
**Why RED now:** The archive prefix contains `repo-release/2.x/`, so one stripped component leaves the suite nested below `2.x/`.
**GREEN when:** The requested slash branch is still selected exactly and the extracted suite—including `VERSION` and manifest inputs—lands directly at the root.
**Validation:** User confirmed by requesting capture of every accepted validation finding.

## Builder Guidance

Certainty level: Firm. Reuse the existing default one-component prefix rather than adding branch-name escaping.

## Full Context

See `do-work/user-requests/UR-083/input.md` for the complete capture provenance.

---
*Source: accepted Finding 23 from the validated external feedback.*

## Triage

**Route: A** - Simple

**Reasoning:** Prefix creation is already centralized in one constant and the missing work is a single end-to-end regression over the existing local Git fixture.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Root Cause

The historical Git fallback built its archive prefix from `repo-<branch>/`. A slash in the branch therefore created an extra path component even though installers always strip exactly one component. Commit `2c82ef12` had already decoupled the prefix from the branch before this accepted finding entered the queue, but it lacked the required production-style regression.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch_test.go` (modified)

**What was done:** Added an end-to-end local fixture that creates `release/2.x`, forces the HTTP route to fail, fetches that exact branch through Git, extracts the resulting archive with production's one-component strip, and asserts the branch-specific `VERSION` at the extraction root. The existing constant `upstream/` prefix remains the production behavior.

## Qualification

Passed — 1 implementation file verified, 4 requirements traced, P-A-U confirmed. The test exercises the public fetch path and real `git archive`/`tar` behavior; missing-branch and derivation coverage remain intact. The unrelated knowledge-base batch was excluded from evidence and staging.

## Testing

**Tests run:** focused archive-fetch tests; `go vet ./... && go test ./... -count=1` in the CLI module; `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All passing, including queue-board ordinary/strict JavaScript lanes and canonical maintainer verification. The optional browser lane was unavailable and skipped by the gate; this backend archive contract has no browser acceptance condition.

**Red-green validation:**
- `TestGitFallbackSlashBranchExtractsAtRootAfterOneComponentStrip`: ✗ under a controlled replay of the historical `upstream-release/2.x/` prefix (`VERSION` absent at extraction root) → ✓ with the retained constant `upstream/` prefix

**New tests added:**
- Exact `release/2.x` clone selection through failed-HTTP fallback, archive generation, `--strip-components=1` extraction, and root manifest assertion

*Verified by work action*

## Review

**Overall: 99%** | 2026-09-01T22:24:26Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — branch text controls only exact clone selection, every Git fallback keeps one prefix component, and production-style extraction places the selected branch's manifest at root.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Orientation

Git fallback archives now have an end-to-end contract proving slash-containing branch names cannot change extraction depth. The do-work CLI prime remains current.
