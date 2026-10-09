## Review: REQ-651

**Approve** — the report gives the maintainer an accurate split of the activity correlation code. Every line citation I checked resolves at bd4d754c. One count in the evidence prose is off by one.
Route A | 5e996eed..bd4d754c (builder 5ac41434, integrator citation fix 7a04e2ec)

### What's built
- One new bundle `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/` (index.html plus three evidence files). It separates the deterministic git reads from the scaffold that exists because a builder commit may lack its `[REQ-NNN]` prefix. It counts five git calls per refresh, cites the measured cost, and ends with the stamp-then-delete versus keep decision, with no recommendation.

### Decisions / risks for you
- None. The report changes no code. The decision it supports belongs to the maintainer.

### Findings

**Important:**
- None.

**Minor:**
- F1. `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html:325` says "Thirteen have the shape `[work run] REQ-NNN step N`". `evidence/prefix-coverage.txt` lists 12 commits with that shape (all REQ-594 steps). The 13th `[work run]` commit, cfcf63e4, reads `[work run] REQ-554 remediation: …` and has no "step". The line should read: "Twelve have the shape `[work run] REQ-NNN step N`, and one more starts with `[work run] REQ-NNN`, each with the id present but not bracketed." — impact-negligible → report only

**Nit:**
- F2. D-01 deletes `activity_correlation.go:111-113` (part of the doc comment for `correlateCommitsToRequests`). Line 110 would then end mid-sentence ("or membership in"). A real deletion would rewrite line 110 instead, so the 54-line count does not change. — impact-negligible → report only

### Requirements Checklist

- [x] ai-report run for REQ-632, one bundle under `ai-reports/<slug>/` with no overwrite. Every path in the range is an addition. — delivered
- [x] Two parts titled exactly "Deterministic git structure" (5 rows) and "Scaffold around agent behaviour" (6 rows). Each row gives file + line range, what it reads, and its origin (REQ-284, UR-135, REQ-632, REQ-636, REQ-646). — delivered
- [x] Git calls per refresh (five, each with a verify.go or activity_correlation.go line) and measured cost (REQ-632 archived lines 71/116/224/256 plus `evidence/git-log-timing.txt`). The figures match the evidence files. — delivered
- [x] Final decision with line counts on each side (303 + 337 kept vs 357 + 429) and "This report does not recommend either option". — delivered
- [x] No change under `skills/`. Render step done by the builder (Playwright, light/dark, 390px). — delivered
- [x] UR-143 F7 maintainer concern ("scaffold around non-deterministic behaviour") answered directly by the scaffold table and the decision section. — delivered

### Acceptance Testing

**Result: Pass** (source review of the citations, no render re-run)
- `git diff --stat 5e996eed..bd4d754c`: 4 files, all under the one ai-reports bundle.
- activity_correlation.go (357 lines) and its test (429 lines) are unchanged between 2eb14357 and bd4d754c. I checked every cited range at bd4d754c: :19-22, :60, :62-64, :70, :96, :111-113, :116-120, :129-136, :139-164 (143-151 guard), :174-191, :209-210, :214-216, :252-303, :274-281, :329-334. Test :141-232 (two tests plus trailing blank), :141-144, :379-406. All match the report's descriptions.
- D-01 deletion sum: 1+1+3+5+26+18 = 54. 357-54 = 303. 429-92 = 337. Correct.
- verify.go at bd4d754c: :1349 rev-parse --abbrev-ref, :1357 rev-parse HEAD, :1397-1442 shared read, :1425-1426 / :1425-1439 for-each-ref, :1448 worktree list, :1470 branch --list, :989 / :1110 / :1381 / :1489 per-leftover direct git. All correct after the D-08 shift of 47. No pre-shift numbers remain in index.html.
- Other citations resolve: commit.md:126 (`[REQ-NNN] {REQ title}`), work-reference.md:416 (worktree naming rule), UR-135 input.md:73-79, :82-84 ("lossy"), :85 (`req-nnnn-*`), :89-91, prime-do-kanban.md:30 (`git-history-evidence`), serve.go:160-174, generate.go:679-686, REQ-632 archived lines 71 (50-90 ms, 72 commits), 116 (74 commits, 8 merges, ~30 ms), 224 (105-129 / 148-157 ms), 256 (2,901 commits, 336 merges, ~0.2 s + 39 ms).
- D-04 corrections are right. The scaffold is 54 lines. "Lossy" is at 82-84. ownedTipInstantsById is read only in activity_correlation.go, so for-each-ref feeds only activity, and verify consumes the other three shared commands.
- `evidence/git-call-count.txt` shows calls=5 for 0, 1 and 4 branches. Timing figures (44-47 / 166 ms; 152-161 / 253 ms) match `git-log-timing.txt`. The 44 / 37 / 34-of-250 / 2026-09-06 figures match `prefix-coverage.txt`, except the F1 count.

### Suggested Additional Testing

- None.

### Scores (on the record — not the headline)

**Overall: 98%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All five requirements delivered |
| Code Quality | 95% | Report accuracy, one off-by-one count (F1) |
| Test Adequacy | N/A | Report-only REQ |
| Scope | 100% | Only the ai-reports bundle |
| Risk | None | — |
| Acceptance | Pass | Citations and figures verified at bd4d754c |

### Follow-ups created
- None (2 findings report only)

## Review

**Overall: 98%** | <INTEGRATOR-TIMESTAMP>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1. index.html:325 says "Thirteen have the shape `[work run] REQ-NNN step N`", but prefix-coverage.txt lists 12 such commits. The 13th `[work run]` commit, cfcf63e4, is a "remediation" commit with no step. — impact-negligible → report only
- F2 (nit). D-01's deletion of activity_correlation.go:111-113 leaves line 110's doc comment mid-sentence. A real deletion rewrites line 110, and the count stays 54. — impact-negligible → report only
**Acceptance:** Pass — every line range, line count, git-call count and cost figure resolves at bd4d754c, except the F1 count
**Restatement sweep:** nothing redefined
**Suggested testing:** 0 items
**Follow-ups created:** None (2 findings report only)

*Reviewed by review-work action*
